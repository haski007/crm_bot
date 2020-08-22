package statistics

import (
	"fmt"
	"time"

	"../../betypes"
	"../../database"
	"../../utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func InitEveryDayStatistics(bot *tgbotapi.BotAPI) {
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
		utils.SendInfoToAdmins(bot, getDailyStatistics())
	}

}

func getDailyStatistics() string {

	var products []betypes.Product

	fromDate := utils.GetTodayStartTime()

	database.ProductsCollection.Find(nil).All(&products)

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