package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"

	apiclient "github.com/eryalito/vigo-bus-core-go-client"
	"github.com/eryalito/vigo-bus-tg-bot/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StopInfoHandler(identity *apiclient.ApiIdentity, bot *tgbotapi.BotAPI, message *tgbotapi.Message, client *apiclient.APIClient) {
	// Parse the message to get the stop id
	stopString := message.Text[1:]
	sendStopInfo(identity, bot, stopString, message.Chat.ID, client)

}
func StopInfoCallbackHandler(identity *apiclient.ApiIdentity, bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, client *apiclient.APIClient) {
	// Parse the callback to get the stop id
	stopString := callback.Data[1:]
	sendStopInfo(identity, bot, stopString, callback.From.ID, client)
}

func sendStopInfo(identity *apiclient.ApiIdentity, bot *tgbotapi.BotAPI, stopString string, chatID int64, client *apiclient.APIClient) {
	stopNumber, err := strconv.Atoi(stopString)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Invalid stop number")
		bot.Send(msg)
		return
	}

	requestStop := apiclient.ApiApiStopsStopNumberGetRequest(client.BusAPI.ApiStopsStopNumberScheduleGet(context.Background(), int32(stopNumber)))
	stopInfo, _, err := requestStop.Execute()
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(chatID, "Stop not found")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(chatID, "Parada "+stopString+": "+*stopInfo.Name)
	bot.Send(msg)

	locationMessage := tgbotapi.NewLocation(chatID, float64(*stopInfo.Location.Lat), float64(*stopInfo.Location.Lon))
	bot.Send(locationMessage)

	requestSchedule := apiclient.ApiApiStopsStopNumberScheduleGetRequest(client.BusAPI.ApiStopsStopNumberScheduleGet(context.Background(), int32(stopNumber)))
	schedule, _, err := requestSchedule.Execute()
	if err != nil {
		log.Println("Error getting stop schedule:", err)
		msg := tgbotapi.NewMessage(chatID, "Schedule not found")
		bot.Send(msg)
		return
	}

	text := "Horario de la parada " + stopString + ":\n\n"
	for _, scheduleItem := range schedule.Schedules {
		text += *scheduleItem.Line.Name + " | " + *scheduleItem.Route + " | " + strconv.FormatInt(int64(*scheduleItem.Time), 10) + "min\n"
	}

	messages := utils.SplitLongMessage(text)
	for i, m := range messages {
		msg := tgbotapi.NewMessage(chatID, m)
		if i == len(messages)-1 {
			buttonRows := make([][]tgbotapi.InlineKeyboardButton, 0)
			favText := "Añadir a favoritos"
			if utils.IsStopInFavorites(identity, int32(stopNumber)) {
				favText = "Borrar de favoritos"
			}
			buttonRows = append(buttonRows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(favText, fmt.Sprintf("/fav %d", stopNumber)),
			))
			buttonRows = append(buttonRows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Recargar", fmt.Sprintf("/%d", stopNumber)),
			))
			buttonRows = addButtonsForLines(buttonRows, schedule.Schedules, stopNumber)
			inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(buttonRows...)
			msg.ReplyMarkup = inlineKeyboard

		}
		bot.Send(msg)
	}
}

func addButtonsForLines(buttons [][]tgbotapi.InlineKeyboardButton, schedules []apiclient.ApiSchedule, stopNumber int) [][]tgbotapi.InlineKeyboardButton {
	addedLines := []int32{}
	for _, scheduleItem := range schedules {
		buttonAlreadyAdded := false
		for _, line := range addedLines {
			if line == *scheduleItem.Line.Id {
				buttonAlreadyAdded = true
				break
			}
		}
		if buttonAlreadyAdded {
			continue
		}
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(*scheduleItem.Line.Name, fmt.Sprintf("/follow %d %d", stopNumber, *scheduleItem.Line.Id)),
		))
		addedLines = append(addedLines, *scheduleItem.Line.Id)
	}
	return buttons
}
