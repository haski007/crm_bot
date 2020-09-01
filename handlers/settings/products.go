package settings

import (
	"encoding/json"
	"strconv"
	"strings"

	"../../betypes"
	"../../database"
	"../../keyboards"
	"../../emoji"

	"github.com/globalsign/mgo/bson"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type m bson.M

// Queue of users who are trying to add new product
var (
	AddProductQueue = make(map[int]*betypes.Product)
)

// GetAllProductsHandler prints all produtcs from "products" collection.
func GetAllProductsHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	var products []betypes.Product

	database.ProductsCollection.Find(bson.M{}).Select(m{"purchases": 0}).Sort("type").All(&products)

	for i, prod := range products {
		prod.Name = "*" + prod.Name + "*"

		j, _ := json.Marshal(prod)
		answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, strconv.Itoa(i)+") "+string(j))
		answer.ParseMode = "Markdown"
		prodKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Remove "+emoji.Minus, "remove_product "+prod.ID.Hex()),
			),
		)
		answer.ReplyMarkup = prodKeyboard
		bot.Send(answer)
	}

	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Here you come!")
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Send(answer)
}

// RemoveProductHandler removes product from "products" collection
func RemoveProductHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	prodID := strings.Split(update.CallbackQuery.Data, " ")[1]

	err := database.ProductsCollection.Remove(bson.M{"_id": bson.ObjectIdHex(prodID)})
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Removing has been FAILED! {"+err.Error()+"}"))
		return
	}


	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "The product has been removed succesfully!")
	answer.ReplyMarkup = keyboards.MainMenu

	bot.Send(answer)
}

// AddProductHandler adds product to database collection "products"
func AddProductHandler(bot *tgbotapi.BotAPI,update tgbotapi.Update) {
	AddProductQueue[update.CallbackQuery.From.ID] = new(betypes.Product)

	bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Enter product name:"))
	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
}

// addProduct prompt user to get name, type and prise of product. Save it in DB
func AddProduct(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.From.ID

	prod := AddProductQueue[update.Message.From.ID]

	var err error
	if prod.Name == "" {
		prod.Name = update.Message.Text

		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose product type:")
		
		typesKeyboard := keyboards.GetTypesKeyboard("protyp")
		typesKeyboard.InlineKeyboard = append(typesKeyboard.InlineKeyboard,
			[]tgbotapi.InlineKeyboardButton{keyboards.MainMenuButton})
		answer.ReplyMarkup = typesKeyboard

		bot.Send(answer)
	} else if prod.PrimeCost == 0.0 {
		prod.PrimeCost, err = strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong type format! Try again"))
			return
		}

		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Enter *selling price* price:")
		answer.ParseMode = "MarkDown"
		bot.Send(answer)
	} else {
		prod.Price, err = strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong type format! Try again"))
			return
		}

		var answer tgbotapi.MessageConfig

		err = database.ProductsCollection.Insert(prod)
		if err != nil {
			answer = tgbotapi.NewMessage(update.Message.Chat.ID, "Product has not beed added {"+err.Error()+"}")
		} else {
			answer = tgbotapi.NewMessage(update.Message.Chat.ID, "Product has been added succesfully!")
		}

		delete(AddProductQueue, userID)
		answer.ReplyMarkup = keyboards.MainMenu
		bot.Send(answer)
	}
}

func AddTypeToProduct(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	t := strings.Join(strings.Fields(update.CallbackQuery.Data)[1:], " ")

	AddProductQueue[update.CallbackQuery.From.ID].Type = t
	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		"Enter prime cost")
	bot.Send(answer)

	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
}