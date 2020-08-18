package main

import (
	"time"
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

	// fromDate := getTodayStarTime()



	// if err := P.Pipe(query).All(&purchases); err != nil {
	// 	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "*FAILED!* {"+err.Error()+"}")
	// 	answer.ParseMode = "MarkDown"
	// 	answer.ReplyMarkup = mainKeyboard
	// 	return answer
	// }

	
	// for i, pur := range purchases {
	// 	fmt.Println(i, ") ", pur)
	// }
	// os.Exit(-1)

	var message string
	// for i, pur := range purchases {
	// 	message += fmt.Sprintf("%s%3d)*%v*\n", emoji.PurchasesDelimiter, i, pur)
	// }
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