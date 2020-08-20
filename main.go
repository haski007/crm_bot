package main

import (
	"log"
	"strings"

	"./botlogs"
	"./emoji"
	"./models"
	"github.com/Haski007/go-errors"
	"github.com/globalsign/mgo/bson"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var removePurchaseQueue = make(map[int]bool)

var makePurchaseQueue = make(map[int]bson.ObjectId)

var addProductQueue = make(map[int]*models.Product)

var mainMenuButton = tgbotapi.NewInlineKeyboardButtonData("........."+emoji.House+"......."+emoji.Tree+
"..Main Menu........"+emoji.HouseWithGarden+"..."+emoji.Car+"....", "home")

var mainKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Purchase "+emoji.Dollar, "purchase"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Configuration "+emoji.Gear, "configs"),
		tgbotapi.NewInlineKeyboardButtonData("Statistics "+emoji.GraphicIncrease, "stats"),
	),
)

func main() {
	err := initMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		errors.Println(err)
		return
	}
	bot.Debug = true

	go initEveryDayStatistics(bot)
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		errors.Println(err)
		return
	}

	logger := botlogs.NewLogger("")

	var resp tgbotapi.MessageConfig

	for update := range updates {
		// ---> Handle keyboard signals.
		if update.CallbackQuery != nil {
			// ---> Validate user
			if !isUser(update.CallbackQuery.From) {
				resp := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "*FORBIDDEN!* you are not registered!\n"+
					"You can register by /register")
				resp.ParseMode = "Markdown"
				bot.Send(resp)
				continue
			}
			switch update.CallbackQuery.Data {
			case "home":
				resp = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
					"........."+emoji.House+"......."+emoji.Tree+"..Main Menu........"+
						emoji.HouseWithGarden+"..."+emoji.Car+"....")
				resp.ReplyMarkup = mainKeyboard
			case "configs":
				resp = ConfigsHandler(update)
			case "add_product":
				resp = AddProductHandler(update)
			case "get_all_products":
				resp = GetAllProductsHandler(bot, update)
				resp.ReplyMarkup = mainKeyboard
			case "purchase":
				resp = GetProductTypesHandler(bot, update)
			case "stats":
				resp = GetStatisticsHandler(bot, update)
			case "curr_day_history":
				resp = GetCurrentDayHistoryHandler(bot, update)
			case "curr_day_stats":
				resp = GetCurrentDayStatsHandler(update)
			case "remove_purchase":
				resp = RemovePurchaseHandler(update)
			case "/test":
				resp = TestHandler(bot, update)
			}

			// Handle callbacks with info
			if strings.Contains(update.CallbackQuery.Data, "remove_product") {
				resp = RemoveProductHandler(update)
				resp.ReplyMarkup = mainKeyboard
			} else if strings.Contains(update.CallbackQuery.Data, "purtyp") {
				resp = GetProductsByTypeHandler(bot, update)
			} else if strings.Contains(update.CallbackQuery.Data, "purname") {
				resp = MakePurchaseHandler(update)
			}

			bot.Send(resp)
		}

		// ---> Handle messages
		if update.Message != nil {

			err := logger.MessageLog(update)
			if err != nil {
				errors.Println(err)
			}

			if command := update.Message.CommandWithAt(); command != "" {
				switch command {
				case "help":
					resp = HelpHandler(update)
				case "register":
					resp = RegisterUser(bot, update, updates)
				case "menu":
					resp = MenuHandler(update)
				case "start":
					resp = StartHandler(update)
				case "users":
					resp = GetAllUsers(update)
				case "remove_user":
					resp = RemoveUserHandler(update)
				default:
					resp = tgbotapi.NewMessage(update.Message.Chat.ID, emoji.Warning+" Unknown command! "+emoji.Warning)				}
			} else {
				if _, ok := addProductQueue[update.Message.From.ID]; ok {
					resp = addProduct(update)
				} else if removePurchaseQueue[update.Message.From.ID] == true {
					resp = removePurchase(update)
				} else if _ , ok := makePurchaseQueue[update.Message.From.ID]; ok {
					resp = makePurchase(update)
				} else {
					resp = tgbotapi.NewMessage(update.Message.Chat.ID,
						emoji.Warning+" It's not a command! "+emoji.Warning)
				}
			}
			bot.Send(resp)
		}
	}
}
