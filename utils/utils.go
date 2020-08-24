package utils

import (
	"time"

	"../betypes"
	"../database"

	"github.com/globalsign/mgo/bson"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var Location, _ = time.LoadLocation("Europe/Kiev")

type m bson.M

func SendInfoToAdmins(bot *tgbotapi.BotAPI, message string) {
	var admins []betypes.User

	err := database.UsersCollection.Find(m{"status": "admin"}).All(&admins)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(370649141, "ALARM: Something went wrong!!!!"))
	}

	for _, user := range admins {
		answer := tgbotapi.NewMessage(int64(user.UserID), message)
		answer.ParseMode = "MarkDown"
		bot.Send(answer)
	}
}

func GetTodayStartTime() time.Time {

	t := time.Now().In(Location)

	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 4, 0, t.Location())
}

func MakeEmojiRow(emoji string, len int) string {
	var row string

	for i := 0; i < len; i++ {
		row += emoji
	}
	return row
}