package handlers

import (
	"context"
	"strconv"
	"strings"

	apiclient "github.com/eryalito/vigo-bus-core-go-client"
	"github.com/eryalito/vigo-bus-tg-bot/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func FavoriteHandler(identity *apiclient.ApiIdentity, bot *tgbotapi.BotAPI, message *tgbotapi.Message, client *apiclient.APIClient) {
	handleFavorite(identity, bot, message.Text, message.Chat.ID, client)
}

func FavoriteCallbackHandler(identity *apiclient.ApiIdentity, bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, client *apiclient.APIClient) {
	handleFavorite(identity, bot, callback.Data, callback.Message.Chat.ID, client)
}

func handleFavorite(identity *apiclient.ApiIdentity, bot *tgbotapi.BotAPI, message string, chatID int64, client *apiclient.APIClient) {
	// Strip the command from the message
	text := strings.TrimSpace(strings.Replace(message, "/fav", "", 1))

	// Check if the user wants to just list the favorites
	if text == "" {
		listFavorites(identity, bot, chatID)
		return
	}

	stopNumber, err := strconv.Atoi(text)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Invalid stop number")
		bot.Send(msg)
		return
	}

	// Get the stop info
	requestStop := apiclient.ApiApiStopsStopNumberGetRequest(client.BusAPI.ApiStopsStopNumberScheduleGet(nil, int32(stopNumber)))
	_, _, err = requestStop.Execute()
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Parada no encontrada")
		bot.Send(msg)
		return
	}

	// Check if the stop is already in the favorites
	if identity.FavoriteStops == nil {
		identity.FavoriteStops = make([]apiclient.ApiStop, 0)
	}
	isFav := utils.IsStopInFavorites(identity, int32(stopNumber))

	if isFav {
		// Remove the stop from the favorites
		request := client.IdentityAPI.ApiUsersProviderUuidFavoriteStopsStopNumberDelete(context.Background(), "telegram", *identity.Uuid, int32(stopNumber))
		_, _, err = client.IdentityAPI.ApiUsersProviderUuidFavoriteStopsStopNumberDeleteExecute(request)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Error al eliminar la parada de favoritos")
			bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(chatID, "Parada eliminada de favoritos")
		bot.Send(msg)
		return
	}

	// Add the stop to the favorites
	request := client.IdentityAPI.ApiUsersProviderUuidFavoriteStopsStopNumberPost(context.Background(), "telegram", *identity.Uuid, int32(stopNumber))
	_, _, err = client.IdentityAPI.ApiUsersProviderUuidFavoriteStopsStopNumberPostExecute(request)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Error al añadir la parada a favoritos")
		bot.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(chatID, "Parada añadida a favoritos")
	bot.Send(msg)
}

func listFavorites(identity *apiclient.ApiIdentity, bot *tgbotapi.BotAPI, chatID int64) {
	if identity.FavoriteStops == nil || len(identity.FavoriteStops) == 0 {
		msg := tgbotapi.NewMessage(chatID, "No tienes paradas favoritas")
		bot.Send(msg)
		return
	}

	text := "Paradas favoritas:\n\n"
	for _, stop := range identity.FavoriteStops {
		text += strconv.Itoa(int(*stop.StopNumber)) + " - " + *stop.Name + "\n"
	}

	for _, m := range utils.SplitLongMessage(text) {
		msg := tgbotapi.NewMessage(chatID, m)
		bot.Send(msg)
	}
}
