package handlers

import (
	apiclient "github.com/eryalito/vigo-bus-core-go-client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartHandler(identity *apiclient.ApiIdentity, bot *tgbotapi.BotAPI, message *tgbotapi.Message, client *apiclient.APIClient) {
	msg := tgbotapi.NewMessage(message.Chat.ID, `
Hola! Soy el bot de Buses de Vigo (no oficial). Para ver lo que puedo hacer, escribe /help
		`)
	bot.Send(msg)
}
