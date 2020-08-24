package keyboards

import (
	"../emoji"
	"../database"

	"github.com/globalsign/mgo/bson"
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
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Show all product types "+emoji.Info, "get_all_types"),
		tgbotapi.NewInlineKeyboardButtonData("Add new type "+emoji.NewButton, "add_type"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Show all products "+emoji.Box, "get_all_products"),
	),
	tgbotapi.NewInlineKeyboardRow(MainMenuButton),
)

var StatsKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Month stats "+emoji.UpLeftArrow, "month_stats"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Today's history"+emoji.UpLeftArrow, "curr_day_history"),
		tgbotapi.NewInlineKeyboardButtonData("Today's stats"+emoji.GraphicIncrease, "curr_day_stats"),
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

var StoreKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("New supply " + emoji.Registered, "supply"),
		tgbotapi.NewInlineKeyboardButtonData("Check storage " + emoji.QuestionMark, "check_storage"),
	),
	tgbotapi.NewInlineKeyboardRow(MainMenuButton),
)

var TypesListKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Remove type "+emoji.Basket, "remove_type"),
	),
	tgbotapi.NewInlineKeyboardRow(
		MainMenuButton,
	),
)


func GetTypesKeyboard(data string) tgbotapi.InlineKeyboardMarkup {

	var types []string

	database.ProductTypesCollection.Find(bson.M{}).Distinct("type", &types)

	countRows := len(types) / 3
	if countRows % 3 != 0 || countRows == 0{
		countRows++
	}
	rows := make([][]tgbotapi.InlineKeyboardButton, countRows)
	var x int
	for i, t := range types {
		if i%3 == 0 && i != 0 {
			x++
		}
		rows[x] = append(rows[x], tgbotapi.NewInlineKeyboardButtonData(t, data + " " + t))
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return keyboard
}