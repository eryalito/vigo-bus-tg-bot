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
	lat := float32(message.Location.Latitude)
	lon := float32(message.Location.Longitude)
	radius := float32(1000)

	request := &apiclient.ApiApiStopsFindLocationImageGetRequest{
		ApiService: client.BusAPI,
	}

	req := request.Lat(lat).Lon(lon).Radius(radius).Limit(9)

	nearbyStops, _, err := req.Execute()
	if err != nil {
		log.Println("Error getting location image:", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "No se ha podido obtener la imagen")
		bot.Send(msg)
		getStopListByLocation(lat, lon, radius, message, bot, client)
		return
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

	sendStopsMessage(message, bot, nearbyStops.Stops)
}

func getStopListByLocation(lat, lon, radius float32, message *tgbotapi.Message, bot *tgbotapi.BotAPI, client *apiclient.APIClient) {
	request := &apiclient.ApiApiStopsFindLocationGetRequest{
		ApiService: client.BusAPI,
	}

	req := request.Lon(lon).Lat(lat).Radius(radius)

	nearbyStops, _, err := req.Execute()
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Failed to get stops")
		bot.Send(msg)
		return
	}

	sendStopsMessage(message, bot, nearbyStops)
}

func sendStopsMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI, stops []apiclient.ApiStop) {
	stopsText := "Paradas cercanas:\n"

	// Draw markers for each stop
	for index, stop := range stops {
		if index > 9 {
			break
		}

		stopsText += fmt.Sprintf("%d: /%d - %s\n", index+1, *stop.StopNumber, *stop.Name)
	}

	for _, m := range utils.SplitLongMessage(stopsText) {
		msg := tgbotapi.NewMessage(message.Chat.ID, m)
		bot.Send(msg)
	}
}
