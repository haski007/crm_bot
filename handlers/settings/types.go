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
		"Введите новый тип продукта:")
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
	answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Новый тип продукта был создан! " + emoji.Check)
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
		message += fmt.Sprintf("---------------------------------------\nТип продукта #%d: *%s*\n\n%v\n", i + 1, t.Type, t.ID)
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

	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Введите тип продукта который вы хотите удалить:")
	bot.Send(answer)

	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
}

func RemoveType(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	defer func(){
		if err := recover(); err != nil {
			message := fmt.Sprintf("ERROR %s: {%v}", emoji.Warning, err)
			answer := tgbotapi.NewMessage(update.Message.Chat.ID, message)
			answer.ReplyMarkup = keyboards.MainMenu
			bot.Send(answer)
		}
	}()
	typeName := update.Message.Text



	if err := database.ProductTypesCollection.Remove(m{"type":typeName}); err != nil {
		delete(RemoveTypeQueue, update.Message.From.ID)
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR "+emoji.Warning+": {"+err.Error()+"}")
		answer.ReplyMarkup = keyboards.MainMenu
		bot.Send(answer)
		return
	}

	if _, err := database.ProductsCollection.RemoveAll(m{"type":typeName}); err != nil {
		delete(RemoveTypeQueue, update.Message.From.ID)
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR "+emoji.Warning+": {"+err.Error()+"}")
		answer.ReplyMarkup = keyboards.MainMenu
		bot.Send(answer)
		return
	}

	delete(RemoveTypeQueue, update.Message.From.ID)
	answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Тип продукта был успешно удалён!" + emoji.Check)
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Send(answer)
}