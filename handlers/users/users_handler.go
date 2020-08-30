package users

import (
	"fmt"
	"strconv"

	"../../betypes"
	"../../database"
	"../../keyboards"
	"../../utils"
	"github.com/globalsign/mgo/bson"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type m bson.M

var (
	RegisterUserQueue = make(map[int]int)
)

func RegisterUserHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if count, _ := database.UsersCollection.Find(bson.M{"user_id": update.Message.From.ID}).Count(); count > 0 {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are registered already!\nUse /menu"))
		return
	}
	
	RegisterUserQueue[update.Message.From.ID] = 3
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Send me a secret password:"))
}

func RegisterUser(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	
	if update.Message.Text != betypes.SECRET_VASSAL_PASSWORD && update.Message.Text != betypes.SECRET_LORD_PASSWORD {	
		if RegisterUserQueue[update.Message.From.ID] == 1 {
			delete(RegisterUserQueue, update.Message.From.ID)
	
			answer := tgbotapi.NewMessage(update.Message.Chat.ID, "You have used 3 tries, try again /register\n"+
				"If you don't have secret password - write to GOD - @pdemian !")
			bot.Send(answer)
			return
		} else {
			RegisterUserQueue[update.Message.From.ID]--
			answer := tgbotapi.NewMessage(update.Message.Chat.ID, "*WRONG PASSWORD!*\nTry again:")
			answer.ParseMode = "MarkDown"
			bot.Send(answer)
			return
		}
	}
	
	var user betypes.User
		
	user.FirstName = update.Message.From.FirstName
	user.LastName = update.Message.From.LastName
	user.UserName = update.Message.From.UserName
	user.UserID = update.Message.From.ID
	if update.Message.Text == betypes.SECRET_LORD_PASSWORD {
		user.Status = "admin"
	} else {
		user.Status = "user"
	}

	go utils.SendInfoToAdmins(bot, fmt.Sprintf("New user has been registred: %s (%s)\nAs *%s*",
		user.FirstName, user.UserName, user.Status))

	
	err := database.UsersCollection.Insert(user)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Registration has been FAILED {"+err.Error()+"}"))
		return
	}
	
	
	delete(RegisterUserQueue, update.Message.From.ID)
	answer := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, "+user.FirstName+":)\n"+
		"You have been succesfully registered with id: "+strconv.Itoa(user.UserID)+"\n"+
		"Start using bot by /menu")
	answer.ReplyMarkup = keyboards.MainMenu
	bot.Send(answer)
}

func IsAdmin(from *tgbotapi.User) bool {
	var user betypes.User
	
	err := database.UsersCollection.Find(m{"user_id": from.ID}).One(&user)
	if err != nil {
		return false
	}

	if user.Status == "admin" {
		return true
	}
	return false
}

func IsUser(from *tgbotapi.User) bool {
	count, err := database.UsersCollection.Find(m{"user_id": from.ID}).Count()
	if err != nil || count < 1 {
		return false
	}

	return true
}
