package handlers

import (
	"../keyboards"
	"../emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func MainMenuHandler(bot *tgbotapi.BotAPI,update tgbotapi.Update) {
	message := "........."+emoji.House+"......."+emoji.Tree+"..\n   *Main Menu*   \n........"+
		emoji.HouseWithGarden+"..."+emoji.Car+"...."
	answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		message, keyboards.MainMenu)
	answer.ParseMode = "MarkDown"
	bot.Send(answer)
}