package main

import (
	"fmt"
	// "./emoji"
	"time"

	"./emoji"
	"./models"
	"github.com/globalsign/mgo/bson"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type m bson.M

func GetPurchasesHistoryHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) tgbotapi.MessageConfig {
	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Purchases history")

	var historyKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("The all day", "curr_day_history"),
			tgbotapi.NewInlineKeyboardButtonData("The Last 5 purchases", "edit_purchases"),
		),
	)
	answer.ReplyMarkup = historyKeyboard

	return answer
}

func GetCurrentDayHistoryHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) tgbotapi.MessageConfig {

	// var product models.Product

	
	var products []models.Product
	fromDate := getTodayStarTime()
	
	query := m{
		"purchases": m{
			"$elemMatch": m{
				"sale_date" : m{
					"$gt" :fromDate,
				},
			},
		},
	}
	
	err := ProductsCollection.Find(query).All(&products)
	if err != nil {
		return tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "ERROR: " + err.Error())
	}
	

	var message string

	var id int
	for _, prod := range products {
		i := len(prod.Purchases) - 1
		for i > -1 && prod.Purchases[i].SaleDate.After(fromDate) {
			message += fmt.Sprintf("%sPurchase #%d\nProduct: %s\nType: %s\nSold: %.2f\nCash: %.2f\nSale Date: %v\n",
				emoji.PurchasesDelimiter, id, prod.Name, prod.Type, prod.Purchases[i].Amount,
				prod.Purchases[i].Amount * prod.Price, prod.Purchases[i].SaleDate.Format("02.01.2006 15:04:05"))
			id++
			i--
		}
	}
	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, message)
	answer.ParseMode = "MarkDown"
	answer.ReplyMarkup = mainKeyboard
	return answer
}


func EditPurchasesHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "*FAILED!* {}")
}

func getTodayStarTime() time.Time {

	location, _ := time.LoadLocation("Europe/Kiev")
	t := time.Now().In(location)

	year, month, day := t.Date()
    return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}