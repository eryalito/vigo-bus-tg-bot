package handlers

import (
	"strconv"
	"strings"

	apiclient "github.com/eryalito/vigo-bus-core-go-client"
	"github.com/eryalito/vigo-bus-tg-bot/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func FindHandler(identity *apiclient.ApiIdentity, bot *tgbotapi.BotAPI, message *tgbotapi.Message, client *apiclient.APIClient) {
	// Strip the command from the message
	text := strings.TrimSpace(strings.Replace(message.Text, "/search", "", 1))

	if len(text) < 3 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Por favor, introduce al menos 3 caracteres")
		bot.Send(msg)
		return
	}

	request := &apiclient.ApiApiStopsFindGetRequest{
		ApiService: client.BusAPI,
	}
	stops, _, err := request.Text(text).Execute()
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Error al buscar la parada")
		bot.Send(msg)
		return
	}

	sendText := "Resultados de la búsqueda:\n\n"

	for _, stop := range stops {
		sendText += *stop.Name + " - /" + strconv.Itoa(int(*stop.StopNumber)) + "\n"
	}

	for _, m := range utils.SplitLongMessage(sendText) {
		msg := tgbotapi.NewMessage(message.Chat.ID, m)
		bot.Send(msg)
	}

}
