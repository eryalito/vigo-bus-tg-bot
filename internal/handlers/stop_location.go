package handlers

import (
	"context"
	"strconv"

	apiclient "github.com/eryalito/vigo-bus-core-go-client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StopLocationHandler(identity *apiclient.ApiIdentity, bot *tgbotapi.BotAPI, message *tgbotapi.Message, client *apiclient.APIClient) {
	// Parse the message to get the stop id
	stopString := message.Text[1:(len(message.Text) - 1)]
	stopNumber, err := strconv.Atoi(stopString)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid stop number")
		bot.Send(msg)
		return
	}

	requestStop := apiclient.ApiApiStopsStopNumberGetRequest(client.BusAPI.ApiStopsStopNumberScheduleGet(context.Background(), int32(stopNumber)))
	stop, _, err := requestStop.Execute()
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Stop not found")
		bot.Send(msg)
		return
	}

	locationMessage := tgbotapi.NewLocation(message.Chat.ID, float64(*stop.Location.Lat), float64(*stop.Location.Lon))
	bot.Send(locationMessage)

}
