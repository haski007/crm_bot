package statistics

import (
	"fmt"
	"time"

	"../../betypes"
	"../../database"
	"../../keyboards"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var MonthStatsQueue = make(map[int]bool)

func MonthStatisticsHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	MonthStatsQueue[update.CallbackQuery.From.ID] = true
	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		`Enter month and year in format - "08.2020":`)

	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
	bot.Send(answer)
}

func GetMonthStatistics(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	fromDate, err := time.Parse("01.2006", update.Message.Text)
	if err != nil {
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "*WRONG!* {"+err.Error()+"}")
		answer.ReplyMarkup = keyboards.MainMenu
		answer.ParseMode = "MarkDown"
		bot.Send(answer)
		return
	}

	toDate := fromDate.AddDate(0, 1, -1).Add(23 * time.Hour + 59 * time.Minute)

	var products []betypes.Product

	database.ProductsCollection.Find(nil).All(&products)

	var totalSum float64
	var totalMoney float64

	var message string = "   "


	for index, prod := range products {
		amount := 0.0
		profit := 0.0
		i := len(prod.Purchases) - 1
		for i > -1 && prod.Purchases[i].SaleDate.After(fromDate) && prod.Purchases[i].SaleDate.Before(toDate) {
			amount += prod.Purchases[i].Amount
			profit = prod.Purchases[i].Amount * prod.Price - prod.Purchases[i].Amount * prod.PrimeCost
			totalSum += prod.Purchases[i].Amount * prod.Price
			totalMoney += prod.Purchases[i].Amount * prod.PrimeCost
			i--
		}
		message += fmt.Sprintf("%-3d) %s  %-5s(*%v*) profit(*%.2f UAH*)\n", index, prod.Name, "sold", amount, profit)
	}

	message += fmt.Sprintf("Total cash: *%v UAH*\nTotal profit: *%.2f UAH*", totalSum, totalSum - totalMoney)
	delete(MonthStatsQueue, update.Message.From.ID)
	answer := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	answer.ParseMode = "MarkDown"
	answer.ReplyMarkup = keyboards.MainMenu

	bot.Send(answer)
}
