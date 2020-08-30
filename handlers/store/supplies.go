package store

import (
	"strconv"
	"strings"

	"../../betypes"
	"../../database"
	"../../emoji"
	"../../keyboards"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/globalsign/mgo/bson"
)

func GetProductTypesHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	typeChoiceKeyboard := keyboards.GetTypesKeyboard("suptyp")
	typeChoiceKeyboard.InlineKeyboard = append(typeChoiceKeyboard.InlineKeyboard,
		[]tgbotapi.InlineKeyboardButton{keyboards.MainMenuButton})

	answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		"Choose type of product...", typeChoiceKeyboard)
	bot.Send(answer)
}

func GetProductsByTypeHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	t := strings.Join(strings.Split(update.CallbackQuery.Data, " ")[1:], " ")
	var prods []betypes.Product

	database.ProductsCollection.Find(bson.M{"type": t}).All(&prods)

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, prod := range prods {
		rows = append(rows, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(prod.Name, "supname "+prod.ID.Hex()),
		})
	}
	rows = append(rows, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("Main menu "+emoji.House, "home")})

	var productsKeyboard = tgbotapi.NewInlineKeyboardMarkup(rows...)
	answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		"Choose product...",
		productsKeyboard)
	bot.Send(answer)
}

func ReceiveSuppliesHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	getID := strings.Split(update.CallbackQuery.Data, " ")[1]
	productID := bson.ObjectIdHex(getID)

	SupplyQueue[update.CallbackQuery.From.ID] = productID
	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
	
	bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Supply quantity:"))
}

func MakeSupply(bot *tgbotapi.BotAPI,update tgbotapi.Update) {
	supplyValue, err := strconv.ParseFloat(update.Message.Text, 4)
	if err != nil {
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong type format!" + emoji.Warning)
		bot.Send(answer)
	}

	who := m{
		"_id":SupplyQueue[update.Message.From.ID],
	}

	query := m{"$inc": m{
		"in_storage": supplyValue,
	}}
	delete(SupplyQueue, update.Message.From.ID)
	
	if err := database.ProductsCollection.Update(who, query); err != nil {
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, emoji.Warning + "ERROR: {"+err.Error()+"}")
		answer.ReplyMarkup = keyboards.MainMenu
		bot.Send(answer)
	}

	answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Supply was succesfully received! " + emoji.Check)
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Send(answer)
}