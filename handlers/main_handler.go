package handlers

import (
	"../keyboards"
	"../emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func MainMenuHandler(update tgbotapi.Update) tgbotapi.MessageConfig {
	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		"........."+emoji.House+"......."+emoji.Tree+"..Main Menu........"+
			emoji.HouseWithGarden+"..."+emoji.Car+"....")
	answer.ReplyMarkup = keyboards.MainMenu
	return answer
}