package main

import (
	"fmt"
	"time"

	"./emoji"
	"./models"
	"github.com/globalsign/mgo/bson"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var location, _ = time.LoadLocation("Europe/Kiev")

type m bson.M

func GetStatisticsHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) tgbotapi.MessageConfig {
	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Purchases history")

	var historyKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Get today's history"+emoji.UpLeftArrow, "curr_day_history"),
			tgbotapi.NewInlineKeyboardButtonData("Get today's stats"+emoji.GraphicIncrease, "curr_day_stats"),
		),
		tgbotapi.NewInlineKeyboardRow(
			mainMenuButton,
		),
	)
	answer.ReplyMarkup = historyKeyboard

	return answer
}

func GetCurrentDayHistoryHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) tgbotapi.MessageConfig {
	var products []models.Product
	fromDate := getTodayStartTime()

	query := m{
		"purchases": m{
			"$elemMatch": m{
				"sale_date": m{
					"$gt": fromDate,
				},
			},
		},
	}

	err := ProductsCollection.Find(query).All(&products)
	if err != nil {
		return tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "ERROR: "+err.Error())
	}

	var message string

	var id int
	for _, prod := range products {
		i := len(prod.Purchases) - 1
		for i > -1 && prod.Purchases[i].SaleDate.After(fromDate) {
			message += fmt.Sprintf("%sPurchase #%d\nProduct: %s\nType: %s\nSold: %.2f\nCash: %.2f\nSale Date: %v\n%s\n",
				emoji.PurchasesDelimiter, id, prod.Name, prod.Type, prod.Purchases[i].Amount,
				prod.Purchases[i].Amount*prod.Price, prod.Purchases[i].SaleDate.In(location).Format("02.01.2006 15:04:05"), prod.Purchases[i].ID.String())
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
	if !isAdmin(update.CallbackQuery.From) {
		return msgNotEnoughPermissions(update.CallbackQuery.Message)
	}
	message := getDailyStatistics()

	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, message)
	answer.ReplyMarkup = mainKeyboard
	return answer
}

func RemovePurchaseHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update, ch tgbotapi.UpdatesChannel) tgbotapi.MessageConfig {
	defer func() {
		err := recover()
		fmt.Println("\n\n\n\n", err)
	}()
	bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Send me id of purchase you want to remove:"))
	update = <-ch

	purchaseID := bson.ObjectIdHex(update.Message.Text)

	who := m{
		"purchases": m{
			"$elemMatch": m{
				"_id": purchaseID,
			},
		},
	}
	query := m{
		"$pull": m{
			"purchases": m{
				"_id": purchaseID,
			},
		},
	}

	err := ProductsCollection.Update(who, query)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR: {"+err.Error()+"}")
	}

	answer := tgbotapi.NewMessage(update.Message.Chat.ID, "An purchase has been succesfully removed!")
	answer.ReplyMarkup = mainKeyboard

	return answer
}

func getTodayStartTime() time.Time {

	t := time.Now().In(location)

	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 4, 0, t.Location())
}
