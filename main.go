package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"./betypes"
	"./botlogs"
	"./emoji"
	"./handlers"
	"./handlers/cashbox"
	"./handlers/purchases"
	"./handlers/settings"
	"./handlers/statistics"
	"./handlers/store"
	"./handlers/users"
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
				resp := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
					emoji.NoEntry+" *FORBIDDEN!* "+emoji.NoEntry+" you are not registered!\n"+
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
				settings.GetAllProductsHandler(bot, update)
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
				purchases.RemovePurchaseHandler(bot, update)
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
			case "end_day":
				cashbox.EndDayHandler(bot, update)
			}

			// Handle callbacks with info
			if strings.Contains(update.CallbackQuery.Data, "remove_product") {
				settings.RemoveProductHandler(bot, update)
			} else if strings.Contains(update.CallbackQuery.Data, "purtyp") {
				purchases.GetProductsByTypeHandler(bot, update)
			} else if strings.Contains(update.CallbackQuery.Data, "purname") {
				purchases.MakePurchaseHandler(bot, update)
			} else if strings.Contains(update.CallbackQuery.Data, "suptyp") {
				store.GetProductsByTypeHandler(bot, update)
			} else if strings.Contains(update.CallbackQuery.Data, "supname") {
				store.ReceiveSuppliesHandler(bot, update)
			} else if strings.Contains(update.CallbackQuery.Data, "edit_product") {
				settings.EditProductHandler(bot, update)
			} else if strings.Contains(update.CallbackQuery.Data, "getprodstyp") {
				settings.GetAllProducts(bot, update)
			} else if _, ok := settings.AddProductQueue[update.CallbackQuery.From.ID];
						ok && strings.Contains(update.CallbackQuery.Data, "protyp") {
				settings.AddTypeToProduct(bot, update)
			} else if strings.Contains(update.CallbackQuery.Data, "edit ") {
				settings.GetEntityToEdit(bot, update)
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
					handlers.CommandHelpHandler(bot, update)
				case "register":
					users.RegisterUserHandler(bot, update)
				case "menu":
					handlers.CommandMenuHandler(bot, update)
				case "start":
					handlers.CommandStartHandler(bot, update)
				case "users":
					handlers.CommandUsersHandler(bot, update)
				case "remove_user":
					handlers.CommandRemoveUserHandler(bot, update)
				case "remove_today_cash":
					handlers.RemoveTodayCash(bot, update)
				case "alert_all":
					handlers.CommandAlertEverybodyHandler(bot, update)
				case "alert_admins":
					handlers.CommandAlertAdminsHandler(bot, update)
				case "test":
					go test(bot, update)
				default:
					resp = tgbotapi.NewMessage(update.Message.Chat.ID, emoji.Warning+" Unknown command! "+emoji.Warning)
				}
			} else {
				if _, ok := settings.AddProductQueue[update.Message.From.ID]; ok {
					settings.AddProduct(bot, update)
				} else if statistics.MonthStatsQueue[update.Message.From.ID] == true  {
					statistics.GetMonthStatistics(bot, update)
				} else if purchases.RemovePurchaseQueue[update.Message.From.ID] == true {
					purchases.RemovePurchase(bot, update)
				} else if _, ok := purchases.MakePurchaseQueue[update.Message.From.ID]; ok {
					purchases.MakePurchase(bot, update)
				} else if _, ok := store.SupplyQueue[update.Message.From.ID]; ok {
					store.MakeSupply(bot, update)
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
				} else if cashbox.EndDayQueue[update.Message.From.ID] == true {
					cashbox.EndDay(bot, update)
				} else if _, ok := users.RegisterUserQueue[update.Message.From.ID]; ok {
					users.RegisterUser(bot, update)
				} else if _, ok := settings.EditProductNameQueue[update.Message.From.ID]; ok{
					settings.EditProductName(bot, update)
				} else if _, ok := settings.EditProductMarginQueue[update.Message.From.ID]; ok{
					settings.EditProductMargin(bot, update)
				} else if _, ok := settings.EditProductPrimeQueue[update.Message.From.ID]; ok{
					settings.EditProductPrime(bot, update)
				} else if _, ok := settings.EditProductPriceQueue[update.Message.From.ID]; ok{
					settings.EditProductPrice(bot, update)
				} else if _, ok := settings.EditProductUnitQueue[update.Message.From.ID]; ok{
					settings.EditProductUnit(bot, update)
				} else  {
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

func test(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	
	for i := 0; i < 29; i++ {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Waited "+strconv.Itoa(i + 1)+" seconds"))
	}
}