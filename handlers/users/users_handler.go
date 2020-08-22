package users

import (
	"fmt"
	"strconv"

	"../../betypes"
	"../../database"
	"../../utils"
	"github.com/globalsign/mgo/bson"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type m bson.M

func RegisterUser(bot *tgbotapi.BotAPI, update tgbotapi.Update, ch tgbotapi.UpdatesChannel) tgbotapi.MessageConfig {
	if count, _ := database.UsersCollection.Find(bson.M{"user_id": update.Message.From.ID}).Count(); count > 0 {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "You are registered already!\nUse /menu")
	}

	var tries int
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Send me secret password:"))
	for {
		update = <-ch
		if update.Message.Text == betypes.SECRET_VASSAL_PASSWORD || update.Message.Text == betypes.SECRET_LORD_PASSWORD {
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

	utils.SendInfoToAdmins(bot, fmt.Sprintf("New user has been registred: %s (%s)", user.FirstName, user.UserName))

	err := database.UsersCollection.Insert(user)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Registration has been FAILED {"+err.Error()+"}")
	}

	return tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, "+user.FirstName+":)\n"+
		"You have been succesfully registered with id: "+strconv.Itoa(user.UserID)+"\n"+
		"Start using bot by /menu")
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
