package keyboards

import (
	"../emoji"
	"../database"

	"github.com/globalsign/mgo/bson"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var MainMenuButton = tgbotapi.NewInlineKeyboardButtonData("Главное Меню "+emoji.House, "home")

var MainMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Продажа "+emoji.Dollar, "purchase"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Настройки "+emoji.Gear, "configs"),
		tgbotapi.NewInlineKeyboardButtonData("Статистика "+emoji.GraphicIncrease, "stats"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Склад "+emoji.Package, "store"),
		tgbotapi.NewInlineKeyboardButtonData("Касса "+emoji.MoneyFace, "cashbox"),
	),
)


// Create keyboard for configs.
var SettingsKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Создать продукт "+emoji.Plus, "add_product"),
		tgbotapi.NewInlineKeyboardButtonData("Показать все продукты "+emoji.Box, "get_all_products"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить тип "+emoji.NewButton, "add_type"),
		tgbotapi.NewInlineKeyboardButtonData("Показать все типы "+emoji.Info, "get_all_types"),
	),
	tgbotapi.NewInlineKeyboardRow(MainMenuButton),
)

var StatsKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Статистика за месяц "+emoji.UpLeftArrow, "month_stats"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("История покупок"+emoji.UpLeftArrow, "curr_day_history"),
		tgbotapi.NewInlineKeyboardButtonData("Статистика за день"+emoji.GraphicIncrease, "curr_day_stats"),
	),
	tgbotapi.NewInlineKeyboardRow(
		MainMenuButton,
	),
)

var HistoryKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Удалить покупку "+emoji.Basket, "remove_purchase"),
	),
	tgbotapi.NewInlineKeyboardRow(
		MainMenuButton,
	),
)

var StoreKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Поставка продукта " + emoji.Registered, "supply"),
		tgbotapi.NewInlineKeyboardButtonData("Проверить склад " + emoji.QuestionMark, "check_storage"),
	),
	tgbotapi.NewInlineKeyboardRow(MainMenuButton),
)

var TypesListKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Удалить тип "+emoji.Basket, "remove_type"),
	),
	tgbotapi.NewInlineKeyboardRow(
		MainMenuButton,
	),
)

var CashboxKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Пополнить кассу "+emoji.DollarBanknote, "plus_cash"),
		tgbotapi.NewInlineKeyboardButtonData("Взять с кассы "+emoji.MoneyWithWings, "minus_cash"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Все транзакции " + emoji.Receipt, "transactions"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Начать день " + emoji.FileBox, "set_start_cash"),
		tgbotapi.NewInlineKeyboardButtonData("Завершить день " + emoji.EndArrow, "end_day"),
	),
	tgbotapi.NewInlineKeyboardRow(
		MainMenuButton,
	),
)

var EditProductKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Название", "edit prod_name"),
		tgbotapi.NewInlineKeyboardButtonData("Маржа", "edit prod_margin"),
		tgbotapi.NewInlineKeyboardButtonData("Себестоимость", "edit prod_prime"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Цена", "edit prod_price"),
		tgbotapi.NewInlineKeyboardButtonData("Единица", "edit prod_unit"),
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