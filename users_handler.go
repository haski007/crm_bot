package main

import (
	"strconv"

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
		update = <- ch
		if update.Message.Text == SECRET_PASSWORD {
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
	if user.UserName == "pdemian" {
		user.Status = "admin"
	} else {
		user.Status = "user"
	}
	
	
	err := UsersCollection.Insert(user)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Registration has been FAILED {"+err.Error()+"}")
	}

	return tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, "+user.FirstName+":)\n"+
		"You have been succesfully registered with id: "+strconv.Itoa(user.UserID)+"\n"+
		"Start using bot by /menu")
}

func ValidateUser(update tgbotapi.Update) tgbotapi.MessageConfig {
	userID := update.Message.From.ID

	var answer tgbotapi.MessageConfig

	if count, _ := UsersCollection.Find(bson.M{"user_id": userID}).Count(); count > 0 {
		answer = tgbotapi.NewMessage(update.Message.Chat.ID,
			"........."+emoji.House+"......."+emoji.Tree+"..Main Menu........"+
			emoji.HouseWithGarden+"..."+emoji.Car+"....")
		answer.ReplyMarkup = mainKeyboard
	} else {
		answer = tgbotapi.NewMessage(update.Message.Chat.ID, "*FORBIDDEN!* you are not registered!\n"+
			"You can register by /register")
		answer.ParseMode = "Markdown"
	}
	return answer
}