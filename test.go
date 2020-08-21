package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func TestHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) tgbotapi.MessageConfig {

	answer := tgbotapi.NewMessage(update.Message.Chat.ID, "message")
	return answer
}
