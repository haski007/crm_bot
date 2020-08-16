package handlers

import (
	"fmt"

	"../emoji"

	"github.com/globalsign/mgo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// ConfigsHandler handle "Configuration" callback (button)
func ConfigsHandler(update tgbotapi.Update) tgbotapi.MessageConfig {
	resp := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.Text)

	// Create keyboard for configs.
	var configsKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Add product " + emoji.Plus, "add_product"),
		),
	)

	resp.ReplyMarkup = configsKeyboard
	return resp
}

// AddProductHandler adds product to database collection "products"
func AddProductHandler(update tgbotapi.Update, productsCollection *mgo.Collection, ch tgbotapi.UpdatesChannel) tgbotapi.MessageConfig {
	
	 := <-ch
	answer := tgbotapi.NewMessage(resp.Message.Chat.ID, resp.Message.Text)

	return answer
}