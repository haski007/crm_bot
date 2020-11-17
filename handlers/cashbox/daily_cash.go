package cashbox

import (
	"fmt"
	"strconv"
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
	SetStartDailyMoneyQueue = make(map[int]bool)
	EndDayQueue = make(map[int]bool)
)

func SetStartDailyMoneyHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// ---> Check if money are already set for today

	fromDate := utils.GetTodayStartTime()
	query := m{
		"date": m{
			"$gt":fromDate.Add(3 * time.Hour),
		},
	}

	if count, _ := database.DailyCashCollection.Find(query).Count(); count > 0 {
		answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			"На сегодня начальная касса уже была установлена!" + emoji.Warning +
			"\nНо вы можете удалить это значение коммандой - /remove_today_cash",
			keyboards.MainMenu)
		bot.Send(answer)
		return
	}

	// ---> Promp User
	bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		"Сколько денег у вас в кассе в начале дня?"))

	SetStartDailyMoneyQueue[update.CallbackQuery.From.ID] = true	
	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
}

func SetStartDailyMoney(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	money, err := strconv.ParseFloat(update.Message.Text, 64)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			emoji.Warning+" Неверный тип данных! Попробуйте ещё раз!"))
		return
	}
	delete(SetStartDailyMoneyQueue, update.Message.From.ID)

	var dailyCash betypes.DailyCash

	dailyCash.Money = money
	dailyCash.User = fmt.Sprintf("%s %s (@%s)", update.Message.From.FirstName, update.Message.From.LastName,
		update.Message.From.UserName)
	dailyCash.Date = time.Now()

	err = database.DailyCashCollection.Insert(dailyCash)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"ERROR "+emoji.Warning +": {"+err.Error()+"}"))
		return
	}

	answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Успешно установлено! " + emoji.Check)
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Send(answer)
}

func EndDayHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	todaySum := utils.GetTodayAllMoney()

	message := fmt.Sprintf("Сколько денег вы хотите добавить к главной кассе?\n\nУ вас есть *%.2f UAH* в дневной кассе",
		todaySum)
	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, message)
	answer.ParseMode = "MarkDown"
	bot.Send(answer)

	EndDayQueue[update.CallbackQuery.From.ID] = true
	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
}

func EndDay(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	money, err := strconv.ParseFloat(update.Message.Text, 64)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"Неверный тип данных! Попробуйте ещё раз"))
		return
	}
	
	
	totalSum := utils.GetTodayAllMoney()
	
	
	if money > totalSum {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			emoji.Warning+" Вы не можете указать больше денег чем в дневной кассе! Попробуйте ещё раз:"))
		return
	} else if money < 0 {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			emoji.Warning+" Отрицательного количества денег не существует!\nВы в школе учились?\n"+
			"Попробуйте ещё раз у вас всё получиться!"))
		return
	}
	delete(EndDayQueue, update.Message.From.ID)

	var transaction betypes.Transaction


	transaction.Comment = "Стандартное Завершение дня!"
	transaction.Diff = money
	transaction.Author = fmt.Sprintf("%s %s (@%s)",
		update.Message.From.FirstName, update.Message.From.LastName, update.Message.From.UserName)
	transaction.DataTime = update.Message.Time().In(utils.Location)
	transaction.ID = bson.NewObjectId()


	
	if err := database.MakeTransaction(&transaction); err != nil {
		answer := tgbotapi.NewMessage(update.Message.Chat.ID,
			"Error "+emoji.Warning+": {"+err.Error()+"}")
		answer.ReplyMarkup = keyboards.MainMenu
		bot.Send(answer)
		return
	}

	answer := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Вы успешно завершили день! " + emoji.Check)
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Send(answer)
}