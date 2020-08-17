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

	rows := make([][]tgbotapi.InlineKeyboardButton, len(types)/3+1)
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

	purchase.ProductID = bson.ObjectIdHex((strings.Split(update.CallbackQuery.Data, " ")[1]))

	var err error
	for {
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Sold amount:"))
		update = <- ch
		purchase.Amount, err = strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong type format! Try again"))
		} else {
			break
		}
	}

	purchase.SaleDate = time.Now()

	err = PurchasesCollection.Insert(purchase)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Purchase has been FAILED!{"+err.Error()+"}")
	}
	return tgbotapi.NewMessage(update.Message.Chat.ID, "Purchase has been added succesfully")
}
