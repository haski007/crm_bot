package statistics

import (
	"fmt"

	"../../betypes"
	"../../database"
	"../../emoji"
	"../../keyboards"
	"../../utils"
	"../users"
	"github.com/globalsign/mgo/bson"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)


type m bson.M

func GetStatisticsHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) tgbotapi.MessageConfig {
	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Purchases history")

	var historyKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Get today's history"+emoji.UpLeftArrow, "curr_day_history"),
			tgbotapi.NewInlineKeyboardButtonData("Get today's stats"+emoji.GraphicIncrease, "curr_day_stats"),
		),
		tgbotapi.NewInlineKeyboardRow(
			keyboards.MainMenuButton,
		),
	)
	answer.ReplyMarkup = historyKeyboard

	return answer
}

func GetCurrentDayHistoryHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) tgbotapi.MessageConfig {
	var products []betypes.Product
	fromDate := utils.GetTodayStartTime()

	query := m{
		"purchases": m{
			"$elemMatch": m{
				"sale_date": m{
					"$gt": fromDate,
				},
			},
		},
	}

	err := database.ProductsCollection.Find(query).All(&products)
	if err != nil {
		return tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "ERROR: "+err.Error())
	}

	var message string

	var id int
	for _, prod := range products {
		i := len(prod.Purchases) - 1
		for i > -1 && prod.Purchases[i].SaleDate.After(fromDate) {
			message += fmt.Sprintf("%sPurchase #%d\nProduct: %s\nType: %s\nSold: %.2f\nCash: %.2f\nSeller: %s\nSale Date: %v\n%s\n",
				emoji.PurchasesDelimiter, id, prod.Name, prod.Type, prod.Purchases[i].Amount,
				prod.Purchases[i].Amount*prod.Price,
				prod.Purchases[i].Seller,
				prod.Purchases[i].SaleDate.In(utils.Location).Format("02.01.2006 15:04:05"),
				prod.Purchases[i].ID.String())
			id++
			i--
		}
	}
	var answer tgbotapi.MessageConfig
	if message != "" {
		answer = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, message)
		answer.ParseMode = "MarkDown"
	} else {
		answer = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "There aren't purchases today yet!")
	}

	var historyKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Remove purchase "+emoji.Basket, "remove_purchase"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("........."+emoji.House+"......."+emoji.Tree+"..Main Menu........"+
				emoji.HouseWithGarden+"..."+emoji.Car+"....", "home"),
		),
	)
	answer.ReplyMarkup = historyKeyboard
	return answer
}

func GetCurrentDayStatsHandler(update tgbotapi.Update) tgbotapi.MessageConfig {
	if !users.IsAdmin(update.CallbackQuery.From) {
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "You have not enough permissions!")
		answer.ReplyMarkup = keyboards.MainMenu
		return answer
	}
	message := getDailyStatistics()

	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, message)
	answer.ReplyMarkup = keyboards.MainMenu
	return answer
}