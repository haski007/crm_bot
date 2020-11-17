package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"../betypes"
	"../database"
	"../emoji"
	"../keyboards"
	"../utils"
	"./users"
	"github.com/globalsign/mgo/bson"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)


type m bson.M

func CommandMenuHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if !users.IsUser(update.Message.From) {
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, emoji.NoEntry+" *FORBIDDEN!* "+emoji.NoEntry+" you are not registered!\n"+
			"You can register by /register")
		answer.ParseMode = "Markdown"
		bot.Send(answer)
		return
	}

	deleteAllQueues(update.Message.From.ID)

	answer := tgbotapi.NewMessage(update.Message.Chat.ID,
		"........."+emoji.House+"......."+emoji.Tree+"..Главное меню........"+
			emoji.HouseWithGarden+"..."+emoji.Car+"....")
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Send(answer)
}

// CommandHelpHandler handle command "/help"
func CommandHelpHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	resp := tgbotapi.NewMessage(update.Message.Chat.ID,
		"This bot was created to help with your small bussines logging/management.\n/menu to start using this.")
	bot.Send(resp)
}

func CommandStartHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	answer :=  tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, "+update.Message.From.FirstName+".\n"+
		"Here is an awesome telegram bot, it can help you to become more involved"+
		"in your small bussines.\n"+
		"To start using bot you need to be registered (/register).\n"+
		"Author: @pdemian\n")
	bot.Send(answer)
}

func CommandUsersHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if !users.IsAdmin(update.Message.From) {
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "У вас недостаточно прав!")
		answer.ReplyMarkup = keyboards.MainMenu
		bot.Send(answer)
		return
	}

	var users []betypes.User

	err := database.UsersCollection.Find(nil).All(&users)
	if err != nil {
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR: {"+err.Error()+"}")
		bot.Send(answer)
		return
	}

	var message string
	for i, user := range users {
		message += fmt.Sprintf("------------------------------------\n"+"Пользователь #%d\n"+
			"Имя: %s\nФамилия: %s\nНикнейм: @%s\nСтатус пользователя: %s\nID пользователя: %d\n%v\n", i+1,
			user.FirstName, user.LastName, user.UserName, user.Status, user.UserID, user.ID)
	}

	answer := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	bot.Send(answer)
}

func CommandRemoveUserHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if !users.IsAdmin(update.Message.From) {
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "У вас недостаточно прав!")
		answer.ReplyMarkup = keyboards.MainMenu
		bot.Send(answer)
		return
	}

	args := strings.Fields(update.Message.CommandArguments())

	for _, arg := range args {
		err := database.UsersCollection.RemoveId(bson.ObjectIdHex(arg))
		if err != nil {
			answer := tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR: {"+err.Error()+"}")
			answer.ReplyMarkup = keyboards.MainMenu
			bot.Send(answer)
			return
		}
	}
	answer := tgbotapi.NewMessage(update.Message.Chat.ID,
		strconv.Itoa(len(args))+" пользователь был удалён!"+emoji.Recycling+"\n")
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Send(answer)
}

func RemoveTodayCash(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	
	fromDate := utils.GetTodayStartTime()
	query := m{
		"date": m{
			"$gt":fromDate,
		},
	}

	database.DailyCashCollection.RemoveAll(query)

	answer := tgbotapi.NewMessage(update.Message.Chat.ID,
		"All clear")
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Send(answer)
}

func CommandAlertEverybodyHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if !users.IsAdmin(update.Message.From) {
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "У вас недостаточно прав! " + emoji.NoEntry)
		answer.ReplyMarkup = keyboards.MainMenu
		bot.Send(answer)
		return
	}

	message := update.Message.CommandArguments()

	utils.SendInfoToUsers(bot, message)
}

func CommandAlertAdminsHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if !users.IsAdmin(update.Message.From) {
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "У вас недостаточно прав! " + emoji.NoEntry)
		answer.ReplyMarkup = keyboards.MainMenu
		bot.Send(answer)
		return
	}

	message := update.Message.CommandArguments()

	utils.SendInfoToAdmins(bot, message)
}