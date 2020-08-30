package purchases

import (
	"fmt"
	"strconv"
	"strings"

	"time"

	"../../betypes"
	"../../database"
	"../../emoji"
	"../../keyboards"
	"../../utils"
	"github.com/globalsign/mgo/bson"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	MakePurchaseQueue = make(map[int]bson.ObjectId)
	RemovePurchaseQueue = make(map[int]bool)
)


type m bson.M

func GetProductTypesHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	typeChoiceKeyboard := keyboards.GetTypesKeyboard("purtyp")
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
			tgbotapi.NewInlineKeyboardButtonData(prod.Name, "purname "+prod.ID.Hex()),
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

func MakePurchaseHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	getID := strings.Split(update.CallbackQuery.Data, " ")[1]
	productID := bson.ObjectIdHex(getID)
	
	MakePurchaseQueue[update.CallbackQuery.From.ID] = productID
	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
	
	bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Sold amount:"))
}

func MakePurchase(bot *tgbotapi.BotAPI,update tgbotapi.Update) {
	var purchase betypes.Purchase

	var err error

	purchase.Amount, err = strconv.ParseFloat(update.Message.Text, 64)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong type format! Try again"))
		return
	}

	purchase.SaleDate = time.Now()
	purchase.ID = bson.NewObjectId()
	purchase.Seller = fmt.Sprintf("%s (@%s)", update.Message.From.FirstName, update.Message.From.UserName)


	// ---> Build query
	who := m{"_id": MakePurchaseQueue[update.Message.From.ID]}
	pushToArray := m{
		"$push": m{
			"purchases": purchase},
		"$inc": m{
				"in_storage": -purchase.Amount,
			},	
	}

	// ---> Add purchase to db
	err = database.ProductsCollection.Update(who, pushToArray)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Purchase has been FAILED!{"+err.Error()+"}"))
		return
	}


	var prod betypes.Product
	var message string

	// ---> Check product in stock
	err = database.ProductsCollection.Find(who).Select(m{"name":1, "in_storage":1}).One(&prod)
	if err != nil {
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Purchase has been FAILED!{"+err.Error()+"}")
		answer.ReplyMarkup = keyboards.MainMenu
		bot.Send(answer)
		return
	} else if prod.InStorage < 0 {
		prod.InStorage = 0
		updateQuery := m{
			"$set":m{
				"in_storage": 0,
			},
		}
		err := database.ProductsCollection.Update(who, updateQuery)
		if err != nil {
			answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Purchase has been FAILED!{"+err.Error()+"}")
			answer.ReplyMarkup = keyboards.MainMenu
			bot.Send(answer)
			return
		}
		message += fmt.Sprintf("*WARNING %s:* %v units of %s left on stock!\n",
			emoji.Warning, prod.InStorage, prod.Name)
		
		go utils.SendInfoToAdmins(bot, message)
	} else if prod.InStorage < 10.0 {
		message += fmt.Sprintf("*WARNING %s:* %v units of %s left on stock!\n",
			emoji.Warning, prod.InStorage, prod.Name)
		go utils.SendInfoToAdmins(bot, message)
	} 
	message += "Purchase has been added succesfully " + emoji.Check

	delete(MakePurchaseQueue, update.Message.From.ID)
	answer := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	answer.ReplyMarkup = keyboards.MainMenu
	answer.ParseMode = "MarkDown"
	bot.Send(answer)

}

func RemovePurchaseHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	RemovePurchaseQueue[update.CallbackQuery.From.ID] = true
	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Send me id of purchase you want to remove:")

	bot.Send(answer)
}

func RemovePurchase(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	purchaseID := bson.ObjectIdHex(update.Message.Text)

	who := m{
		"purchases": m{
			"$elemMatch": m{
				"_id": purchaseID,
			},
		},
	}

	// ---> Adding wrong amount to the stock
	var purchases []betypes.Purchase


	if err := database.ProductsCollection.Find(who).Distinct("purchases", &purchases); err != nil {
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, emoji.Warning + "ERROR: {"+err.Error()+"}")
		answer.ReplyMarkup = keyboards.MainMenu
		bot.Send(answer)
		return
	}

	var quantity float64

	for _, pur := range purchases {
		if pur.ID == purchaseID {
			quantity = pur.Amount
			break
		}
	}

	query := m{
		"$pull": m{
			"purchases": m{
				"_id": purchaseID,
			},
		},
		"$inc": m{
			"in_storage": quantity,
		},
	}

	

	// ---> Remove wrong purchase from db
	err := database.ProductsCollection.Update(who, query)
	if err != nil {
		delete(RemovePurchaseQueue, update.Message.From.ID)
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR: {"+err.Error()+"}"))
		return
	}

	delete(RemovePurchaseQueue, update.Message.From.ID)
	answer := tgbotapi.NewMessage(update.Message.Chat.ID, "An purchase has been succesfully removed! " + emoji.Check)
	answer.ReplyMarkup = keyboards.MainMenu

	bot.Send(answer)
}
