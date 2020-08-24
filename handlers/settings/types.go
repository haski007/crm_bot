package settings

import (
	"fmt"

	"../../database"
	"../../emoji"
	"../../keyboards"

	"github.com/globalsign/mgo/bson"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	AddTypeQueue = make(map[int]bool)
	RemoveTypeQueue = make(map[int]bool)
)

type prodType struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Type string `bson:"type"`
}

func AddNewTypeHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	AddTypeQueue[update.CallbackQuery.From.ID] = true

	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, 
		"Enter new type:")
	bot.Send(answer)

	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
}

func AddNewType(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	t := update.Message.Text

	err := database.ProductTypesCollection.Insert(prodType{Type:t})
	if err != nil {
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR "+emoji.Warning+": {"+err.Error()+"}")
		answer.ReplyMarkup = keyboards.MainMenu
		bot.Send(answer)
		return
	}

	delete(AddProductQueue, update.Message.From.ID)
	answer := tgbotapi.NewMessage(update.Message.Chat.ID, "New type has been succesfully added! " + emoji.Check)
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Send(answer)
}

func ShowAllProductsHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var types []prodType

	err := database.ProductTypesCollection.Find(bson.M{}).All(&types)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
			"ERROR "+emoji.Warning+": {"+err.Error()+"}"))
		return
	}

	var message string
	for i, t := range types {
		message += fmt.Sprintf("---------------------------------------\nType #%d: *%s*\n\n%v\n", i + 1, t.Type, t.ID)
	}

	answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		message,
		keyboards.TypesListKeyboard)
	answer.ParseMode = "MarkDown"
	bot.Send(answer)
}

func RemoveTypeHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	RemoveTypeQueue[update.CallbackQuery.From.ID] = true

	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Send me id of type you want to remove:")
	bot.Send(answer)

	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
}

func RemoveType(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	typeID := bson.ObjectIdHex(update.Message.Text)

	if err := database.ProductTypesCollection.RemoveId(typeID); err != nil {
		delete(RemoveTypeQueue, update.Message.From.ID)
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR "+emoji.Warning+": {"+err.Error()+"}")
		answer.ReplyMarkup = keyboards.MainMenu
		bot.Send(answer)
		return
	}

	answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Type has been succesfully removed!" + emoji.Check)
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Send(answer)
}