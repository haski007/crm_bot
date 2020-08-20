package main

import (
	"fmt"
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

	countRows := len(types) / 3
	if countRows == 0 {
		countRows++
	}
	rows := make([][]tgbotapi.InlineKeyboardButton, countRows)
	var x int
	for i, t := range types {
		if i%3 == 0 && i != 0 {
			x++
		}
		rows[x] = append(rows[x], tgbotapi.NewInlineKeyboardButtonData(t, "purtyp "+t))
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
			tgbotapi.NewInlineKeyboardButtonData(prod.Name, "purname "+prod.ID.Hex()),
		})
	}
	rows = append(rows, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("Main menu "+emoji.House, "home")})

	var productsKeyboard = tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Choose product...")
	msg.ReplyMarkup = productsKeyboard

	return msg
}

func MakePurchaseHandler(update tgbotapi.Update) tgbotapi.MessageConfig {
	
	getID := strings.Split(update.CallbackQuery.Data, " ")[1]
	productID := bson.ObjectIdHex(getID)
	
	makePurchaseQueue[update.CallbackQuery.From.ID] = productID
	
	return tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Sold amount:")
}

func makePurchase(update tgbotapi.Update) tgbotapi.MessageConfig {
	var purchase models.Purchase

	var err error

	purchase.Amount, err = strconv.ParseFloat(update.Message.Text, 64)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong type format! Try again")
	}

	purchase.SaleDate = time.Now()
	fmt.Println(purchase.SaleDate.Format("02.01.2006 15:04:05"))
	purchase.ID = bson.NewObjectId()

	// ---> Build query
	who := m{"_id": makePurchaseQueue[update.Message.From.ID]}
	pushToArray := m{"$push": m{"purchases": purchase}}
	err = ProductsCollection.Update(who, pushToArray)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Purchase has been FAILED!{"+err.Error()+"}")
	}

	delete(makePurchaseQueue, update.Message.From.ID)	
	return tgbotapi.NewMessage(update.Message.Chat.ID, "Purchase has been added succesfully")

}

func RemovePurchaseHandler(update tgbotapi.Update) tgbotapi.MessageConfig {
	removePurchaseQueue[update.CallbackQuery.From.ID] = true
	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Send me id of purchase you want to remove:")

	return answer
}

func removePurchase(update tgbotapi.Update) tgbotapi.MessageConfig {
	purchaseID := bson.ObjectIdHex(update.Message.Text)

	who := m{
		"purchases": m{
			"$elemMatch": m{
				"_id": purchaseID,
			},
		},
	}

	query := m{
		"$pull": m{
			"purchases": m{
				"_id": purchaseID,
			},
		},
	}

	err := ProductsCollection.Update(who, query)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR: {"+err.Error()+"}")
	}

	delete(removePurchaseQueue, update.Message.From.ID)
	answer := tgbotapi.NewMessage(update.Message.Chat.ID, "An purchase has been succesfully removed!")
	answer.ReplyMarkup = mainKeyboard

	return answer
}