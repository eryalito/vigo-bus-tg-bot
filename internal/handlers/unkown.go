package handlers

import (
	apiclient "github.com/eryalito/vigo-bus-core-go-client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func UnknownHandler(identity *apiclient.ApiIdentity, bot *tgbotapi.BotAPI, message *tgbotapi.Message, client *apiclient.APIClient) {
	// Noop
}
