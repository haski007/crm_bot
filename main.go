package main

import (
	"fmt"
	"log"
	"strings"

	"./betypes"
	"./botlogs"
	"./emoji"
	"./keyboards"
	"./handlers"
	"./handlers/settings"
	"./handlers/statistics"
	"./handlers/purchases"
	"./handlers/users"
	"./handlers/store"
	"github.com/Haski007/go-errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(betypes.BOT_TOKEN)
	if err != nil {
		errors.Println(err)
		return
	}
	bot.Debug = true

	go statistics.InitEveryDayStatistics(bot)
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		errors.Println(err)
		return
	}

	logger := botlogs.NewLogger("")

	defer func() {

		if er := recover(); er != nil {
			bot.Send(tgbotapi.NewMessage(370649141, fmt.Sprintf("%+v\n", er)))
		}
	}()

	for update := range updates {
		var resp tgbotapi.MessageConfig

		// ---> Handle keyboard signals.
		if update.CallbackQuery != nil {
			// ---> Validate user
			if !users.IsUser(update.CallbackQuery.From) {
				resp := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "*FORBIDDEN!* you are not registered!\n"+
					"You can register by /register")
				resp.ParseMode = "Markdown"
				bot.Send(resp)
				continue
			}
			switch update.CallbackQuery.Data {
			case "home":
				handlers.MainMenuHandler(bot, update)
			case "configs":
				settings.SettingsHandler(bot, update)
			case "add_product":
				resp = settings.AddProductHandler(update)
			case "get_all_products":
				resp = settings.GetAllProductsHandler(bot, update)
			case "purchase":
				purchases.GetProductTypesHandler(bot, update)
			case "stats":
				statistics.GetStatisticsHandler(bot, update)
			case "curr_day_history":
				statistics.GetCurrentDayHistoryHandler(bot, update)
			case "curr_day_stats":
				statistics.GetCurrentDayStatsHandler(bot, update)
			case "remove_purchase":
				resp = purchases.RemovePurchaseHandler(update)
			case "store":
				store.StoreHandler(bot, update)
			}

			// Handle callbacks with info
			if strings.Contains(update.CallbackQuery.Data, "remove_product") {
				resp = settings.RemoveProductHandler(update)
				resp.ReplyMarkup = keyboards.MainMenu
			} else if strings.Contains(update.CallbackQuery.Data, "purtyp") {
				purchases.GetProductsByTypeHandler(bot, update)
			} else if strings.Contains(update.CallbackQuery.Data, "purname") {
				resp = purchases.MakePurchaseHandler(bot, update)
			}

			if resp.Text != "" {
				bot.Send(resp)
			}
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
					resp = handlers.CommandHelpHandler(update)
				case "register":
					resp = users.RegisterUser(bot, update, updates)
				case "menu":
					resp = handlers.CommandMenuHandler(update)
				case "start":
					resp = handlers.CommandStartHandler(update)
				case "users":
					resp = handlers.CommandUsersHandler(update)
				case "remove_user":
					resp = handlers.CommandRemoveUserHandler(update)
				default:
					resp = tgbotapi.NewMessage(update.Message.Chat.ID, emoji.Warning+" Unknown command! "+emoji.Warning)
				}
			} else {
				if _, ok := settings.AddProductQueue[update.Message.From.ID]; ok {
					resp = settings.AddProduct(update)
				} else if purchases.RemovePurchaseQueue[update.Message.From.ID] == true {
					resp = purchases.RemovePurchase(update)
				} else if _, ok := purchases.MakePurchaseQueue[update.Message.From.ID]; ok {
					resp = purchases.MakePurchase(update)
				} else {
					resp = tgbotapi.NewMessage(update.Message.Chat.ID,
						emoji.Warning+" It's not a command! "+emoji.Warning)
				}
			}
			bot.Send(resp)
		}
	}
}
