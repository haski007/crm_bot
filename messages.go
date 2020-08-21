package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"./keyboards"
)

func msgNotEnoughPermissions(message *tgbotapi.Message) tgbotapi.MessageConfig {
	answer := tgbotapi.NewMessage(message.Chat.ID, "You have not enough permissions!")
	answer.ReplyMarkup = keyboards.MainMenu
	return answer
}
