package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"../betypes"
	"../database"
	"../emoji"
	"../keyboards"
	"./users"
	"github.com/globalsign/mgo/bson"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)


func CommandMenuHandler(update tgbotapi.Update) tgbotapi.MessageConfig {
	if !users.IsUser(update.Message.From) {
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "*FORBIDDEN!* you are not registered!\n"+
			"You can register by /register")
		answer.ParseMode = "Markdown"
		return answer
	}

	answer := tgbotapi.NewMessage(update.Message.Chat.ID,
		"........."+emoji.House+"......."+emoji.Tree+"..Main Menu........"+
			emoji.HouseWithGarden+"..."+emoji.Car+"....")
	answer.ReplyMarkup = keyboards.MainMenu
	return answer
}

// CommandHelpHandler handle command "/help"
func CommandHelpHandler(update tgbotapi.Update) tgbotapi.MessageConfig {
	resp := tgbotapi.NewMessage(update.Message.Chat.ID,
		"This bot was created to help with your small bussines logging/management.\n/menu to start using this.")
	return resp
}

func CommandStartHandler(update tgbotapi.Update) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, "+update.Message.From.FirstName+".\n"+
		"Here is an awesome telegram bot, it can help you to become more involved"+
		"in your small bussines.\n"+
		"To start using bot you need to be registered (/register).\n"+
		"Author: @pdemian\n")
}

func CommandUsersHandler(update tgbotapi.Update) tgbotapi.MessageConfig {
	if !users.IsAdmin(update.Message.From) {
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "You have not enough permissions!")
		answer.ReplyMarkup = keyboards.MainMenu
		return answer
	}

	var users []betypes.User

	err := database.UsersCollection.Find(nil).All(&users)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR: {"+err.Error()+"}")
	}

	var message string
	for i, user := range users {
		message += fmt.Sprintf("------------------------------------\n"+"User #%d\n"+
			"First Name: %s\nLast Name: %s\nUsername: @%s\nUser status: %s\nUser id: %d\n%v\n", i+1,
			user.FirstName, user.LastName, user.UserName, user.Status, user.UserID, user.ID)
	}

	answer := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	return answer
}

func CommandRemoveUserHandler(update tgbotapi.Update) tgbotapi.MessageConfig {
	if !users.IsAdmin(update.Message.From) {
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "You have not enough permissions!")
		answer.ReplyMarkup = keyboards.MainMenu
		return answer
	}

	args := strings.Fields(update.Message.CommandArguments())

	for _, arg := range args {
		err := database.UsersCollection.RemoveId(bson.ObjectIdHex(arg))
		if err != nil {
			answer := tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR: {"+err.Error()+"}")
			answer.ReplyMarkup = keyboards.MainMenu
			return answer
		}
	}
	answer := tgbotapi.NewMessage(update.Message.Chat.ID,
		strconv.Itoa(len(args))+" users has been removed!"+emoji.Recycling+"\n")
	answer.ReplyMarkup = keyboards.MainMenu
	return answer
}