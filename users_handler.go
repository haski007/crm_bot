package main

import (
	"fmt"
	"strconv"
	"strings"

	"./emoji"
	"./models"
	"github.com/globalsign/mgo/bson"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func RegisterUser(bot *tgbotapi.BotAPI, update tgbotapi.Update, ch tgbotapi.UpdatesChannel) tgbotapi.MessageConfig {

	if count, _ := UsersCollection.Find(bson.M{"user_id": update.Message.From.ID}).Count(); count > 0 {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "You are registered already!\nUse /menu")
	}

	var tries int
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Send me secret password:"))
	for {
		update = <-ch
		if update.Message.Text == SECRET_VASSAL_PASSWORD || update.Message.Text == SECRET_LORD_PASSWORD {
			break
		} else if tries > 1 {
			return tgbotapi.NewMessage(update.Message.Chat.ID, "You have used 3 tries, try again /register\n"+
				"If you don't have secret password - write to GOD - @pdemian !")
		} else {
			answer := tgbotapi.NewMessage(update.Message.Chat.ID, "*WRONG PASSWORD!*\nTry again")
			answer.ParseMode = "Markdown"
			bot.Send(answer)
			tries++
		}
	}

	var user models.User

	user.FirstName = update.Message.From.FirstName
	user.LastName = update.Message.From.LastName
	user.UserName = update.Message.From.UserName
	user.UserID = update.Message.From.ID
	if update.Message.Text == SECRET_LORD_PASSWORD {
		user.Status = "admin"
	} else {
		user.Status = "user"
	}

	sendInfoToAdmins(bot, fmt.Sprintf("New user has been registred: %s (%s)", user.FirstName, user.UserName))

	err := UsersCollection.Insert(user)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Registration has been FAILED {"+err.Error()+"}")
	}

	return tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, "+user.FirstName+":)\n"+
		"You have been succesfully registered with id: "+strconv.Itoa(user.UserID)+"\n"+
		"Start using bot by /menu")
}

func MenuHandler(update tgbotapi.Update) tgbotapi.MessageConfig {
	if !isUser(update.Message.From) {
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, "*FORBIDDEN!* you are not registered!\n"+
			"You can register by /register")
		answer.ParseMode = "Markdown"
		return answer
	}

	answer := tgbotapi.NewMessage(update.Message.Chat.ID,
		"........."+emoji.House+"......."+emoji.Tree+"..Main Menu........"+
			emoji.HouseWithGarden+"..."+emoji.Car+"....")
	answer.ReplyMarkup = mainKeyboard
	return answer
}

func GetAllUsers(update tgbotapi.Update) tgbotapi.MessageConfig {
	if !isAdmin(update.Message.From) {
		return msgNotEnoughPermissions(update.CallbackQuery.Message)
	}
	
	var users []models.User

	err := UsersCollection.Find(nil).All(&users)
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

func RemoveUserHandler(update tgbotapi.Update) tgbotapi.MessageConfig {
	if !isAdmin(update.Message.From) {
		return msgNotEnoughPermissions(update.CallbackQuery.Message)
	}

	args := strings.Fields(update.Message.CommandArguments())

	for _, arg := range args {
		err := UsersCollection.RemoveId(bson.ObjectIdHex(arg))
		if err != nil {
			answer := tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR: {"+err.Error()+"}")
			answer.ReplyMarkup = mainKeyboard
			return answer
		}
	}
	answer := tgbotapi.NewMessage(update.Message.Chat.ID,
		strconv.Itoa(len(args))+" users has been removed!"+emoji.Recucling+"\n")
	answer.ReplyMarkup = mainKeyboard
	return answer
}

func isAdmin(from *tgbotapi.User) bool {
	var user models.User

	err := UsersCollection.Find(m{"user_id": from.ID}).One(&user)
	if err != nil {
		return false
	}

	if user.Status == "admin" {
		return true
	}
	return false
}

func isUser(from *tgbotapi.User) bool {
	count, err := UsersCollection.Find(m{"user_id": from.ID}).Count()
	if err != nil || count < 1 {
		return false
	}

	return true
}
