package settings

import (
	"../../keyboards"
	"../../emoji"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)



// SettingsHandler handle "Configuration" callback (button)
func SettingsHandler(update tgbotapi.Update) tgbotapi.MessageConfig {
	resp := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		emoji.Gear+" Configurations "+emoji.Gear)

	resp.ReplyMarkup = keyboards.SettingsKeyboard
	return resp
}
