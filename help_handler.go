package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// HelpHandler handle command "/help"
func HelpHandler(update tgbotapi.Update) tgbotapi.MessageConfig {
	resp := tgbotapi.NewMessage(update.Message.Chat.ID,
		"This bot was created to help with your small bussines logging/management.\n/menu to start using this.")
	return resp
}

func StartHandler(update tgbotapi.Update) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, "+update.Message.From.FirstName+".\n"+
		"Here is an awesome telegram bot, it can help you to become more involved"+
		"in your small bussines.\n"+
		"To start using bot you need to be registered (/register).\n"+
		"Author: @pdemian\n")
}
