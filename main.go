package main

import (
	"log"

	"../errors"
	"./botlogs"
	"./emoji"
	"./handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var mainKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Configuration " + emoji.Gear, "configs"),
	),
)

func main() {
	err := initMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		errors.Println(err)
		return
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		errors.Println(err)
		return
	}

	logger := botlogs.NewLogger("")

	var resp tgbotapi.MessageConfig

	for update := range updates {

		// ---> Handle keyboard signals.
		if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case "home":
				resp = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
					"Main menu " + emoji.House)
				resp.ReplyMarkup = mainKeyboard
			case "configs":
				resp = handlers.ConfigsHandler(update)
			case "add_product":
				handlers.AddProductHandler(update, productsCollection, updates)
			}

			bot.Send(resp)
		}

		// ---> Handle messages
		if update.Message != nil {
		
			err := logger.MessageLog(update)
			if err != nil {
				errors.Println(err)
			}
			
			
			switch update.Message.Text {
			case "/help":
				resp = handlers.HelpHandler(update)
				resp.ReplyMarkup = mainKeyboard
			default:
				resp = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				resp.ReplyMarkup = mainKeyboard
			}
			bot.Send(resp)
		}
	}
}