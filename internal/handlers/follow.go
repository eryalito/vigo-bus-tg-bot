package handlers

import (
	"fmt"
	"log"
	"regexp"
	"time"

	apiclient "github.com/eryalito/vigo-bus-core-go-client"
	"github.com/eryalito/vigo-bus-tg-bot/internal/config"
	"github.com/eryalito/vigo-bus-tg-bot/internal/follow"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var TIME_VALUES = []int{3, 5, 10, 15}

func FollowCallbackHandler(identity *apiclient.ApiIdentity, bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, client *apiclient.APIClient) {
	// If the callback is called long after the original message, then fail and tell the user to reload the message

	if callback.Message.Date+config.MaxTimeAllowedForCallbackAfterMessage < int(time.Now().Unix()) {
		if config.Debug {
			log.Println("Callback called after the message has expired", callback.Message.Date, int(time.Now().Unix()))
		}
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "El mensaje ha expirado, por favor recarga el mensaje.")
		bot.Send(msg)
		return
	}

	// If the message follows the format `/follow [stop] [lineID]`, then update the original message with the buttons to select the time
	// If the message follows the format `/follow [stop] [lineID] [time]`, then just generate a new follow on the line on that stop at that time

	regex := regexp.MustCompile(`/follow (\d+) (\d+)$`)
	matches := regex.FindStringSubmatch(callback.Data)
	if len(matches) > 0 {
		// update the message buttons
		var stopNumber, lineID int
		fmt.Sscanf(callback.Data, "/follow %d %d", &stopNumber, &lineID)
		updateMessageButtons(identity, bot, callback.Message, stopNumber, lineID, client)
		return
	}

	regex = regexp.MustCompile(`/follow (\d+) (\d+) (\d+)$`)
	matches = regex.FindStringSubmatch(callback.Data)
	if len(matches) > 0 {
		// create a new follow
		var stopNumber, lineID, time int
		fmt.Sscanf(callback.Data, "/follow %d %d %d", &stopNumber, &lineID, &time)
		createFollow(identity, bot, callback.Message, stopNumber, lineID, time, client)
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Siguiendo la línea")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "No se puede procesar la petición")
	bot.Send(msg)
}

func createFollow(identity *apiclient.ApiIdentity, bot *tgbotapi.BotAPI, message *tgbotapi.Message, stopNumber, lineID, desiredETA int, client *apiclient.APIClient) {
	follow.FollowTaskManagerInstance.AddTask(follow.FollowTask{
		StopNumber: stopNumber,
		LineID:     lineID,
		DesiredETA: desiredETA,
		ChatID:     message.Chat.ID,
		Client:     client,
		Bot:        bot,
		Manager:    follow.FollowTaskManagerInstance,
	}, time.Duration(config.FollowTaskEvalInterval)*time.Second)
}

func updateMessageButtons(identity *apiclient.ApiIdentity, bot *tgbotapi.BotAPI, message *tgbotapi.Message, stopNumber int, lineID int, client *apiclient.APIClient) {
	// Get the stop info
	// Get the line info
	// Get the times for that line on that stop
	log.Println("Updating message buttons")
	// Generate the buttons
	buttonRows := make([][]tgbotapi.InlineKeyboardButton, 0)

	buttonRows = append(buttonRows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Recargar", fmt.Sprintf("/%d", stopNumber)),
	))
	for _, val := range TIME_VALUES {
		text := fmt.Sprintf("/follow %d %d %d", stopNumber, lineID, val)
		buttonRows = append(buttonRows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d min", val), text)))
	}
	// Update the message
	msg := tgbotapi.NewEditMessageReplyMarkup(message.Chat.ID, message.MessageID, tgbotapi.NewInlineKeyboardMarkup(buttonRows...))
	bot.Send(msg)

}
