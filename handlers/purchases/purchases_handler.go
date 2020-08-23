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
	"github.com/globalsign/mgo/bson"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	MakePurchaseQueue = make(map[int]bson.ObjectId)
	RemovePurchaseQueue = make(map[int]bool)
)


type m bson.M

func GetProductTypesHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	var types []string

	database.ProductsCollection.Find(bson.M{}).Distinct("type", &types)

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

func MakePurchaseHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) tgbotapi.MessageConfig {

	getID := strings.Split(update.CallbackQuery.Data, " ")[1]
	productID := bson.ObjectIdHex(getID)
	
	MakePurchaseQueue[update.CallbackQuery.From.ID] = productID
	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
	
	return tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Sold amount:")
}

func MakePurchase(update tgbotapi.Update) tgbotapi.MessageConfig {
	var purchase betypes.Purchase

	var err error

	purchase.Amount, err = strconv.ParseFloat(update.Message.Text, 64)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong type format! Try again")
	}

	purchase.SaleDate = time.Now()
	purchase.ID = bson.NewObjectId()
	purchase.Seller = fmt.Sprintf("%s (@%s)", update.Message.From.FirstName, update.Message.From.UserName)


	// ---> Build query
	who := m{"_id": MakePurchaseQueue[update.Message.From.ID]}
	pushToArray := m{"$push": m{"purchases": purchase}}
	err = database.ProductsCollection.Update(who, pushToArray)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Purchase has been FAILED!{"+err.Error()+"}")
	}

	delete(MakePurchaseQueue, update.Message.From.ID)
	answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Purchase has been added succesfully")
	answer.ReplyMarkup = keyboards.MainMenu
	return answer

}

func RemovePurchaseHandler(update tgbotapi.Update) tgbotapi.MessageConfig {
	RemovePurchaseQueue[update.CallbackQuery.From.ID] = true
	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Send me id of purchase you want to remove:")

	return answer
}

func RemovePurchase(update tgbotapi.Update) tgbotapi.MessageConfig {
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

	err := database.ProductsCollection.Update(who, query)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR: {"+err.Error()+"}")
	}

	delete(RemovePurchaseQueue, update.Message.From.ID)
	answer := tgbotapi.NewMessage(update.Message.Chat.ID, "An purchase has been succesfully removed!")
	answer.ReplyMarkup = keyboards.MainMenu

	return answer
}
