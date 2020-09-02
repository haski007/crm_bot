package settings

import (
	"strconv"
	"strings"

	"../../database"
	"../../emoji"
	"../../keyboards"
	"github.com/globalsign/mgo/bson"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	EditProductQueue = make(map[int]bson.ObjectId)
	EditProductNameQueue = make(map[int]bson.ObjectId)
	EditProductMarginQueue = make(map[int]bson.ObjectId)
	EditProductPrimeQueue = make(map[int]bson.ObjectId)
	EditProductPriceQueue = make(map[int]bson.ObjectId)
	EditProductUnitQueue = make(map[int]bson.ObjectId)
)

// EditProductHandler edit product from "products" collection
func EditProductHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	prodID := strings.Split(update.CallbackQuery.Data, " ")[1]

	EditProductQueue[update.CallbackQuery.From.ID] = bson.ObjectIdHex(prodID)

	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Choose what do you want to edit!")
	answer.ReplyMarkup = keyboards.EditProductKeyboard

	bot.Send(answer)
}

func GetEntityToEdit(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	entity := strings.Fields(update.CallbackQuery.Data)[1]
	productID := EditProductQueue[update.CallbackQuery.From.ID]

	delete(EditProductQueue, update.CallbackQuery.From.ID)

	switch entity {
	case "prod_name":
		EditProductNameQueue[update.CallbackQuery.From.ID] = productID
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, 
			"Send a new value for name of product!"))
	case "prod_margin":
		EditProductMarginQueue[update.CallbackQuery.From.ID] = productID
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, 
			"Send a new value for margin of product!"))
	case "prod_prime":
		EditProductPrimeQueue[update.CallbackQuery.From.ID] = productID
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, 
			"Send a new value for prime cost of product!"))
	case "prod_price":
		EditProductPriceQueue[update.CallbackQuery.From.ID] = productID
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, 
			"Send a new value for price of product!"))
	case "prod_unit":
		EditProductUnitQueue[update.CallbackQuery.From.ID] = productID
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, 
			"Send a new value for unit of product!"))
	}

	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
}

func EditProductName(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	who := m{
		"_id":EditProductNameQueue[update.Message.From.ID],
	}
	delete(EditProductNameQueue, update.Message.From.ID)

	query := m{
		"$set":m{
			"name":update.Message.Text,
		},
	}

	database.ProductsCollection.Update(who, query)
	answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Product's name was successfully edited " + emoji.Check)
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Send(answer)
}

func EditProductMargin(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	newValue, err := strconv.ParseFloat(update.Message.Text, 64)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong type format! Try again"))
		return
	}

	who := m{
		"_id":EditProductMarginQueue[update.Message.From.ID],
	}
	delete(EditProductMarginQueue, update.Message.From.ID)


	query := m{
		"$set":m{
			"margin": newValue,
		},
	}

	database.ProductsCollection.Update(who, query)
	answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Product's margin was successfully edited " + emoji.Check)
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Send(answer)
}

func EditProductPrime(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	newValue, err := strconv.ParseFloat(update.Message.Text, 64)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong type format! Try again"))
		return
	}

	who := m{
		"_id":EditProductPrimeQueue[update.Message.From.ID],
	}
	delete(EditProductPrimeQueue, update.Message.From.ID)


	query := m{
		"$set":m{
			"prime_cost": newValue,
		},
	}

	database.ProductsCollection.Update(who, query)
	answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Product's prime cost was successfully edited " + emoji.Check)
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Send(answer)
}

func EditProductPrice(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	newValue, err := strconv.ParseFloat(update.Message.Text, 64)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong type format! Try again"))
		return
	}

	who := m{
		"_id":EditProductPriceQueue[update.Message.From.ID],
	}
	delete(EditProductPriceQueue, update.Message.From.ID)


	query := m{
		"$set":m{
			"price": newValue,
		},
	}

	database.ProductsCollection.Update(who, query)
	answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Product's price was successfully edited " + emoji.Check)
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Send(answer)
}

func EditProductUnit(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	who := m{
		"_id":EditProductUnitQueue[update.Message.From.ID],
	}
	delete(EditProductUnitQueue, update.Message.From.ID)

	query := m{
		"$set":m{
			"unit":update.Message.Text,
		},
	}

	database.ProductsCollection.Update(who, query)
	answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Product's unit was successfully edited " + emoji.Check)
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Send(answer)
}