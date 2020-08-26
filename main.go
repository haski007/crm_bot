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
	"./handlers/cashbox"
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
			err := logger.MessageLog(update.CallbackQuery.From, update.CallbackQuery.Data)
			if err != nil {
				errors.Println(err)
			}

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
				settings.AddProductHandler(bot, update)
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
			case "month_stats":
				statistics.MonthStatisticsHandler(bot, update)
			case "remove_purchase":
				resp = purchases.RemovePurchaseHandler(update)
			case "store":
				store.StoreHandler(bot, update)
			case "supply":
				store.GetProductTypesHandler(bot, update)
			case "check_storage":
				store.ShowStorageHandler(bot, update)
			case "add_type":
				settings.AddNewTypeHandler(bot, update)
			case "get_all_types":
				settings.ShowAllProductsHandler(bot, update)
			case "remove_type":
				settings.RemoveTypeHandler(bot, update)
			case "cashbox":
				cashbox.CashboxHandler(bot, update)
			case "plus_cash":
				cashbox.PlusCashHandler(bot, update)
			case "minus_cash":
				cashbox.MinusCashHandler(bot, update)
			case "transactions":
				cashbox.TransactionsHistoryHandler(bot, update)
			case "set_start_cash":
				cashbox.SetStartDailyMoneyHandler(bot, update)
			case "get_start_cash":
				cashbox.GetStartDailyMoneyHandler(bot, update)
			}

			// Handle callbacks with info
			if strings.Contains(update.CallbackQuery.Data, "remove_product") {
				resp = settings.RemoveProductHandler(update)
				resp.ReplyMarkup = keyboards.MainMenu
			} else if strings.Contains(update.CallbackQuery.Data, "purtyp") {
				purchases.GetProductsByTypeHandler(bot, update)
			} else if strings.Contains(update.CallbackQuery.Data, "purname") {
				resp = purchases.MakePurchaseHandler(bot, update)
			} else if strings.Contains(update.CallbackQuery.Data, "suptyp") {
				store.GetProductsByTypeHandler(bot, update)
			} else if strings.Contains(update.CallbackQuery.Data, "supname") {
				resp = store.ReceiveSuppliesHandler(bot, update)
			} else if _, ok := settings.AddProductQueue[update.CallbackQuery.From.ID];
						ok && strings.Contains(update.CallbackQuery.Data, "protyp") {
				settings.AddTypeToProduct(bot, update)
			}

			if resp.Text != "" {
				bot.Send(resp)
			}
		}

		// ---> Handle messages
		if update.Message != nil {
			err := logger.MessageLog(update.Message.From, update.Message.Text)
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
				case "remove_today_cash":
					handlers.RemoveTodayCash(bot, update)
				default:
					resp = tgbotapi.NewMessage(update.Message.Chat.ID, emoji.Warning+" Unknown command! "+emoji.Warning)
				}
			} else {
				if _, ok := settings.AddProductQueue[update.Message.From.ID]; ok {
					resp = settings.AddProduct(update)
				} else if statistics.MonthStatsQueue[update.Message.From.ID] == true  {
					resp = statistics.GetMonthStatistics(update)
				} else if purchases.RemovePurchaseQueue[update.Message.From.ID] == true {
					resp = purchases.RemovePurchase(update)
				} else if _, ok := purchases.MakePurchaseQueue[update.Message.From.ID]; ok {
					resp = purchases.MakePurchase(bot, update)
				} else if _, ok := store.SupplyQueue[update.Message.From.ID]; ok {
					resp = store.MakeSupply(update)
				} else if settings.AddTypeQueue[update.Message.From.ID] == true{
					settings.AddNewType(bot, update)
				} else if settings.RemoveTypeQueue[update.Message.From.ID] == true {
					settings.RemoveType(bot, update)
				} else if _, ok := cashbox.PlusCashQueue[update.Message.From.ID]; ok {
					cashbox.PlusCash(bot, update)
				} else if _, ok := cashbox.MinusCashQueue[update.Message.From.ID]; ok {
					cashbox.MinusCash(bot, update)
				} else if cashbox.TransactionsHostoryQueue[update.Message.From.ID] == true {
					cashbox.ShowTransactionsHistory(bot, update)
				} else if cashbox.SetStartDailyMoneyQueue[update.Message.From.ID] == true {
					cashbox.SetStartDailyMoney(bot, update)
				} else if cashbox.GetStartDailyMoneyQueue[update.Message.From.ID] == true {
					cashbox.GetStartDailyMoney(bot, update)
				} else {
					resp = tgbotapi.NewMessage(update.Message.Chat.ID,
						emoji.Warning+" It's not a command! "+emoji.Warning)
				}
			}
			if resp.Text != "" {
				bot.Send(resp)
			}
		}
	}
}
