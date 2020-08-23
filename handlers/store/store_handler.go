package store

import (
	"fmt"

	"../../betypes"
	"../../database"
	"../../emoji"
	"../../keyboards"
	"../../utils"
	"github.com/globalsign/mgo/bson"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type m bson.M

var SupplyQueue = make(map[int]bson.ObjectId)

func StoreHandler(bot *tgbotapi.BotAPI,update tgbotapi.Update) {
	emojiRow := utils.MakeEmojiRow(emoji.Package,6)

	message := fmt.Sprintf("%s\n\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t*STORE*\n%s\n", emojiRow, emojiRow)

	answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID, message, keyboards.StoreKeyboard)
	answer.ParseMode = "MarkDown"
	bot.Send(answer)
}

func ShowStorageHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var products []betypes.Product

	err := database.ProductsCollection.Find(nil).Select(m{"name":1, "in_storage":1}).All(&products)
	if err != nil {
		answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			emoji.Warning + "ERROR: {"+err.Error()+"}",
			keyboards.MainMenu)
		bot.Send(answer)
	}
	
	var message string
	for i, prod := range products {
		message += fmt.Sprintf("%2d) *%s* - in stock: *%v*\n", i + 1, prod.Name, prod.InStorage)
	}
	answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		message,
		tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(keyboards.MainMenuButton)))
	answer.ParseMode = "MarkDown"
	bot.Send(answer)
}