package handlers

import (
	apiclient "github.com/eryalito/vigo-bus-core-go-client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HelpHandler(identity *apiclient.ApiIdentity, bot *tgbotapi.BotAPI, message *tgbotapi.Message, client *apiclient.APIClient) {
	msg := tgbotapi.NewMessage(message.Chat.ID, `
/start - Mensaje de bienvenida
/help - Mensaje de ayuda
/search [texto] - Buscar paradas por texto
/[número parada] - Ver el horario en tiempo real de la parada
/[número parada]L - Ver la ubicación de la parada

Envíame una ubicación para ver las paradas cercanas!
		`)
	bot.Send(msg)
}
