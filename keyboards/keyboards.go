package keyboards

import (
	"../emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var MainMenuButton = tgbotapi.NewInlineKeyboardButtonData("Main Menu "+emoji.House, "home")

var MainMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Purchase "+emoji.Dollar, "purchase"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Configuration "+emoji.Gear, "configs"),
		tgbotapi.NewInlineKeyboardButtonData("Statistics "+emoji.GraphicIncrease, "stats"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Store "+emoji.Package, "store"),
	),
)


// Create keyboard for configs.
var SettingsKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Add product "+emoji.Plus, "add_product"),
		tgbotapi.NewInlineKeyboardButtonData("Show all products "+emoji.Box, "get_all_products"),
	),
	tgbotapi.NewInlineKeyboardRow(MainMenuButton),
)

var StatsKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Get today's history"+emoji.UpLeftArrow, "curr_day_history"),
		tgbotapi.NewInlineKeyboardButtonData("Get today's stats"+emoji.GraphicIncrease, "curr_day_stats"),
	),
	tgbotapi.NewInlineKeyboardRow(
		MainMenuButton,
	),
)

var HistoryKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Remove purchase "+emoji.Basket, "remove_purchase"),
	),
	tgbotapi.NewInlineKeyboardRow(
		MainMenuButton,
	),
)

// var StoreKeyboard = tgbotapi.NewInlineKeyboardMarkup(
// 	tgbotapi.NewInlineKeyboardRow(
// 		,
// 	),
// 	tgbotapi.NewInlineKeyboardRow(MainMenuButton),
// )