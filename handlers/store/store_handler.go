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

	err := database.ProductsCollection.Find(nil).Select(m{
		"name":1,
		"in_storage":1,
		"type":1,
		"unit":1,
		}).Sort("type").All(&products)
	if err != nil {
		answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			emoji.Warning + "ERROR: {"+err.Error()+"}",
			keyboards.MainMenu)
		bot.Send(answer)
	}
	
	var message string
	var t string
	for i, prod := range products {
		alert := ""
		if prod.InStorage < 10 {
			alert = emoji.RedTrianle
		} else {
			alert = emoji.GreenCircle
		}

		if prod.Type != t {
			t = prod.Type
			message += "*" + emoji.PawPrint + emoji.PawPrint + emoji.PawPrint + emoji.PawPrint +
			t + emoji.PawPrint + emoji.PawPrint + emoji.PawPrint + emoji.PawPrint + "*\n"
		}

		message += fmt.Sprintf("%02d) %s*%s* - in stock: *%v %s*\n",
			i + 1, alert, prod.Name, prod.InStorage, prod.Unit)
	}
	answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		message,
		tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(keyboards.MainMenuButton)))
	answer.ParseMode = "MarkDown"
	bot.Send(answer)
}