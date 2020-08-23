package store

import (
	"fmt"

	"../../emoji"
	"../../keyboards"
	"../../utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func StoreHandler(bot *tgbotapi.BotAPI,update tgbotapi.Update) {
	emojiRow := utils.MakeEmojiRow(emoji.Package,12)

	message := fmt.Sprintf("%s\n\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t*STORE*\n%s\n", emojiRow, emojiRow)

	answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID, message, keyboards.SettingsKeyboard)
	answer.ParseMode = "MarkDown"
	bot.Send(answer)
}