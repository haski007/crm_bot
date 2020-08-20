package main

import (
	"fmt"
	"time"

	"./models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func initEveryDayStatistics(bot *tgbotapi.BotAPI) {
	t := time.Now()
	n := time.Date(t.Year(), t.Month(), t.Day(), 23, 0, 0, 0, t.Location())
	d := n.Sub(t)

	if d < 0 {
		n = n.Add(24 * time.Hour)
		d = n.Sub(t)
	}

	for {
		time.Sleep(d)
		d = (24 * time.Hour)
		sendInfoToAdmins(bot, getDailyStatistics())
	}

}

func getDailyStatistics() string {

	var products []models.Product

	fromDate := getTodayStartTime()

	ProductsCollection.Find(nil).All(&products)

	var totalSum float64

	var message string = "   "

	for index, prod := range products {
		amount := 0.0
		i := len(prod.Purchases) - 1
		for i > -1 && prod.Purchases[i].SaleDate.After(fromDate) {
			amount += prod.Purchases[i].Amount
			totalSum += prod.Purchases[i].Amount * prod.Price
			i--
		}
		message += fmt.Sprintf("%-3d) %-30s %-5s (%.2f)\n", index, prod.Name, "sold", amount)
	}

	message += fmt.Sprintf("Total: %.2f\n", totalSum)

	return message
}

func sendInfoToAdmins(bot *tgbotapi.BotAPI, message string) {
	var admins []models.User

	err := UsersCollection.Find(m{"status": "admin"}).All(&admins)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(370649141, "ALARM: Something went wrong!!!!"))
	}
	for _, user := range admins {
		bot.Send(tgbotapi.NewMessage(int64(user.UserID), message))
	}
}
