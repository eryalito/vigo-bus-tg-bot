package handlers

import (
	"encoding/base64"
	"fmt"
	"log"

	apiclient "github.com/eryalito/vigo-bus-core-go-client"
	"github.com/eryalito/vigo-bus-tg-bot/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func LocationHandler(identity *apiclient.ApiIdentity, bot *tgbotapi.BotAPI, message *tgbotapi.Message, client *apiclient.APIClient) {
	request := &apiclient.ApiApiStopsFindLocationImageGetRequest{
		ApiService: client.BusAPI,
	}
	req := request.Lat(float32(message.Location.Latitude)).Lon(float32(message.Location.Longitude)).Radius(1000).Limit(9)

	nearbyStops, _, err := req.Execute()
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Failed to get stops")
		bot.Send(msg)
		return
	}

	stopsText := "Paradas cercanas:\n"

	// Draw markers for each stop
	for index, stop := range nearbyStops.Stops {
		if index > 9 {
			break
		}

		stopsText += fmt.Sprintf("%d: /%d - %s\n", index+1, *stop.StopNumber, *stop.Name)
	}

	// Get the image bytes from the response
	imgData, err := base64.StdEncoding.DecodeString(*nearbyStops.Image)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Failed to get image")
		bot.Send(msg)
		return
	}

	photo := tgbotapi.NewPhoto(message.Chat.ID, tgbotapi.FileBytes{
		Name:  "map.png",
		Bytes: imgData,
	})
	bot.Send(photo)

	for _, m := range utils.SplitLongMessage(stopsText) {
		msg := tgbotapi.NewMessage(message.Chat.ID, m)
		bot.Send(msg)
	}
}
