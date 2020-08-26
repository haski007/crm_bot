package cashbox

import (
	"fmt"
	"strconv"
	"time"

	"../../betypes"
	"../../emoji"
	"../../keyboards"
	"../../utils"
	"../../database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	SetStartDailyMoneyQueue = make(map[int]bool)
	GetStartDailyMoneyQueue = make(map[int]bool)
)

func SetStartDailyMoneyHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// ---> Check if money are already set for today

	fromDate := utils.GetTodayStartTime()
	query := m{
		"date": m{
			"$gt":fromDate,
		},
	}

	if count, _ := database.DailyCashCollection.Find(query).Count(); count > 0 {
		answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			"Today's cash is already set!" + emoji.Warning +
			"\nBut you can change it by deleting today's start cash /remove_today_cash",
			keyboards.MainMenu)
		bot.Send(answer)
		return
	}

	// ---> Promp User
	bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		"How much money you have in cashbox?"))

	SetStartDailyMoneyQueue[update.CallbackQuery.From.ID] = true	
	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
}

func SetStartDailyMoney(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	money, err := strconv.ParseFloat(update.Message.Text, 64)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"Wrong type format! {"+err.Error()+"}"))
		return
	}
	delete(SetStartDailyMoneyQueue, update.Message.From.ID)

	var dailyCash betypes.DailyCash

	dailyCash.Money = money
	dailyCash.User = fmt.Sprintf("%s %s (@%s)", update.Message.From.FirstName, update.Message.From.LastName,
		update.Message.From.UserName)
	dailyCash.Date = time.Now()

	err = database.DailyCashCollection.Insert(dailyCash)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"ERROR "+emoji.Warning +": {"+err.Error()+"}"))
		return
	}

	answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Succesfully set! " + emoji.Check)
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Send(answer)
}


func GetStartDailyMoneyHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	// ---> Promp User
	bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		"What date you whant to see? (in format '25.12.2020')"))

	GetStartDailyMoneyQueue[update.CallbackQuery.From.ID] = true	
	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
}

func GetStartDailyMoney(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	fromDate, err := time.Parse("02.01.2006", update.Message.Text)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"ERROR "+emoji.Warning +": {"+err.Error()+"}"))
		return
	}
	delete(GetStartDailyMoneyQueue, update.Message.From.ID)

	toDate := fromDate.Add(23 * time.Hour + 59 * time.Minute)
	query := m{
		"date": m{
			"$gt":fromDate,
			"$lt":toDate,
		},
	}

	var dailyCash betypes.DailyCash
	
	if err := database.DailyCashCollection.Find(query).One(&dailyCash); err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		return
	}

	message := fmt.Sprintf("Date: %s\nThe start cash was: %.2f\nUser: %s\n",
		dailyCash.Date.Format("02.01.2006 15:04:05"),
		dailyCash.Money,
		dailyCash.User)

	answer := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Send(answer)
}