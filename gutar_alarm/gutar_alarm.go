package main

import (
	"fmt"
	"net/http"

	"gopkg.in/telegram-bot-api.v4"
)

const (
	BotToken   = "702709092:AAFNpizXtquhDmL1qhVVdYrdDbOuZ7mD8GM"
	WebhookURL = "https://167.99.137.134:88"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		panic(err)
	}

	updates := bot.ListenForWebhook("/")

	go http.ListenAndServe(":88", nil)
	fmt.Println("start listen :88")

	for update := range updates {
		if "gutar" == update.Message.Text {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Yes, this is the gutar",
			))
		} else {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Do you need "+update.Message.Text,
			))
		}
	}
}
