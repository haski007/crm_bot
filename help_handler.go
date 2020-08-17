package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// HelpHandler handle command "/help"
func HelpHandler(update tgbotapi.Update) tgbotapi.MessageConfig {
	resp := tgbotapi.NewMessage(update.Message.Chat.ID,
		"It's bot to help with your small bussines logging/management")

	return resp
}
