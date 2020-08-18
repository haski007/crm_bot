package main

import (
	"strconv"
	"strings"
	"time"

	"./emoji"
	"./models"
	"github.com/globalsign/mgo/bson"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func GetProductTypesHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) tgbotapi.MessageConfig {

	var types []string

	ProductsCollection.Find(bson.M{}).Distinct("type", &types)

	countRows := len(types)/3
	if countRows == 0 {
		countRows++
	}
	rows := make([][]tgbotapi.InlineKeyboardButton, countRows)
	var x int
	for i, t := range types {
		if i%3 == 0 && i != 0 {
			x++
		}
		rows[x] = append(rows[x], tgbotapi.NewInlineKeyboardButtonData(t, "purchase_product_type "+t))
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("Main menu "+emoji.House, "home")})

	typeChoiceKeyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Choose type of product...")
	msg.ReplyMarkup = typeChoiceKeyboard

	return msg
}

func GetProductsByTypeHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) tgbotapi.MessageConfig {
	t := strings.Join(strings.Split(update.CallbackQuery.Data, " ")[1:], " ")
	var prods []models.Product

	ProductsCollection.Find(bson.M{"type": t}).All(&prods)

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, prod := range prods {
		rows = append(rows, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(prod.Name, "purchase_product_name "+prod.ID.Hex()),
		})
	}
	rows = append(rows, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("Main menu "+emoji.House, "home")})

	var productsKeyboard = tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Choose product...")
	msg.ReplyMarkup = productsKeyboard

	return msg
}

func MakePurchaseHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update, ch tgbotapi.UpdatesChannel) tgbotapi.MessageConfig {
	var purchase models.Purchase

	getID := strings.Split(update.CallbackQuery.Data, " ")[1]
	productID := bson.ObjectIdHex(getID)

	var err error
	bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Sold amount:"))
	for {
		update = <- ch
		purchase.Amount, err = strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong type format! Try again"))
		} else {
			break
		}
	}
	
	timezone, err := time.LoadLocation("Europe/Kiev")
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR: {" + err.Error() + "}")
	}
	purchase.SaleDate = time.Now().In(timezone)
	purchase.ID = bson.NewObjectId()
	
	// ---> Build query
	who := m{"_id" : productID}
	pushToArray := m{"$push":m{"purchases":purchase}}
	err = ProductsCollection.Update(who, pushToArray)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Purchase has been FAILED!{"+err.Error()+"}")
	}
	return tgbotapi.NewMessage(update.Message.Chat.ID, "Purchase has been added succesfully")
}
