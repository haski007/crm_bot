package cashbox

import (
	"fmt"
	"math"
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

type m bson.M

var (
	PlusCashQueue = make(map[int]*betypes.Transaction)
	MinusCashQueue = make(map[int]*betypes.Transaction)
	TransactionsHostiryQueue = make(map[int]bool)
)

func CashboxHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	message := fmt.Sprintf("%s\n*Cashbox*\n%s", emoji.MoneyFace, emoji.MoneyFace)

	answer := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		message,
		keyboards.CashboxKeyboard)
	answer.ParseMode = "MarkDown"
	bot.Send(answer)
}

func PlusCashHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	PlusCashQueue[update.CallbackQuery.From.ID] = new(betypes.Transaction)
	message := "How much money do you want to add?"

	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		message)
	bot.Send(answer)

	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
}

func PlusCash(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if PlusCashQueue[update.Message.From.ID].Diff == 0.0 {
		m, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong type format! {"+err.Error()+"}")
			bot.Send(answer)
			return
		}

		PlusCashQueue[update.Message.From.ID].Diff = m
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "What is you purpose?"))
		return
	}
	PlusCashQueue[update.Message.From.ID].Comment = update.Message.Text

	PlusCashQueue[update.Message.From.ID].Author = fmt.Sprintf("%s %s (@%s)",
		update.Message.From.FirstName, update.Message.From.LastName, update.Message.From.UserName)
	PlusCashQueue[update.Message.From.ID].DataTime = update.Message.Time().In(utils.Location)
	PlusCashQueue[update.Message.From.ID].ID = bson.NewObjectId()

	// ---> database manipulations
	who := m{"type": "general"}
	pushToArray := m{
		"$push": m{
			"transactions": PlusCashQueue[update.Message.From.ID]},
		"$inc": m{
				"money": PlusCashQueue[update.Message.From.ID].Diff,
			},	
	}
	err := database.CashboxCollection.Update(who, pushToArray)
	if err != nil {
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR "+emoji.Warning+": {"+err.Error()+"}")
		bot.Send(answer)
		return
	}
	delete(PlusCashQueue, update.Message.From.ID)

	answer := tgbotapi.NewMessage(update.Message.Chat.ID,
		"The transaction was successfully completed! "+  emoji.Check)
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Send(answer)
}

func MinusCashHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	MinusCashQueue[update.CallbackQuery.From.ID] = new(betypes.Transaction)
	message := "How much money did you take?"

	answer := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		message)
	bot.Send(answer)

	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
}

func MinusCash(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if MinusCashQueue[update.Message.From.ID].Diff == 0.0 {
		m, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong type format! {"+err.Error()+"}")
			bot.Send(answer)
			return
		}

		MinusCashQueue[update.Message.From.ID].Diff = -m
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "What is you purpose?"))
		return
	}
	MinusCashQueue[update.Message.From.ID].Comment = update.Message.Text

	MinusCashQueue[update.Message.From.ID].Author = fmt.Sprintf("%s %s (@%s)",
		update.Message.From.FirstName, update.Message.From.LastName, update.Message.From.UserName)
	MinusCashQueue[update.Message.From.ID].DataTime = update.Message.Time().In(utils.Location)
	MinusCashQueue[update.Message.From.ID].ID = bson.NewObjectId()

	// ---> database manipulations
	who := m{"type": "general"}
	pushToArray := m{
		"$push": m{
			"transactions": MinusCashQueue[update.Message.From.ID]},
		"$inc": m{
				"money": MinusCashQueue[update.Message.From.ID].Diff,
			},	
	}
	err := database.CashboxCollection.Update(who, pushToArray)
	if err != nil {
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR "+emoji.Warning+": {"+err.Error()+"}")
		bot.Send(answer)
		return
	}
	delete(MinusCashQueue, update.Message.From.ID)

	answer := tgbotapi.NewMessage(update.Message.Chat.ID,
		"The transaction was successfully completed! "+  emoji.Check)
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Send(answer)
}

func TransactionsHistoryHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	TransactionsHostiryQueue[update.CallbackQuery.From.ID] = true
	bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "How much day you what to see?"))

	bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
}


type transaction struct {
	diff float64
	author string
	comment string
	datatime string
	id bson.ObjectId

}

func ShowTransactionsHistory(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	days, err := strconv.Atoi(update.Message.Text)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, emoji.Warning + " Wrong type format! {"+err.Error()+"}"))
		return
	}

	year, month, day := time.Now().AddDate(0, 0, -days).Date()
	fromDate := time.Date(year, month, day, 0, 0, 0, 0, utils.Location)

	var cashbox betypes.Cashbox

	err = database.CashboxCollection.Find(m{"type":"general"}).One(&cashbox)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR "+ emoji.Warning +": {"+err.Error()+"}"))
		return
	}

	var transactions []transaction
	var message string
	var i = len(cashbox.Transactions) - 1
	for i > -1 && cashbox.Transactions[i].DataTime.After(fromDate) {
		transactions = append(transactions, transaction{
			cashbox.Transactions[i].Diff,
			cashbox.Transactions[i].Author,
			cashbox.Transactions[i].Comment,
			cashbox.Transactions[i].DataTime.Format("02.01.2006 15:04:05"),
			cashbox.Transactions[i].ID,
		})
		
		i--
	}
	
	var index = 1
	var event string
	for i := len(transactions) - 1; i > -1; i-- {
		if transactions[i].diff > 0 {
			message += emoji.GreenDelimiter
			event = "Plus"
		} else {
			event = "Minus"
			message += emoji.RedDelimiter
		}
			
		message += fmt.Sprintf("Transaction #%d\n%s: %.2f\nAuthor: %s\nComment: %s\nDataTime: %s\n%v\n",
			index, event, math.Abs(transactions[i].diff), transactions[i].author,
			transactions[i].comment, transactions[i].datatime,
			transactions[i].id)
		index++
	}
	
	message += fmt.Sprintf("\n%sMoney in Cashbox: *%.2f UAH*", emoji.DollarDelimiter, cashbox.Money)
	answer := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	answer.ReplyMarkup = keyboards.MainMenu
	answer.ParseMode = "MarkDown"
	bot.Send(answer)
}