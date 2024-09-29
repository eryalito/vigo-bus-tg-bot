package handlers

import (
	"context"
	"log"
	"strconv"

	apiclient "github.com/eryalito/vigo-bus-core-go-client"
	"github.com/eryalito/vigo-bus-tg-bot/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StopInfoHandler(identity *apiclient.ApiIdentity, bot *tgbotapi.BotAPI, message *tgbotapi.Message, client *apiclient.APIClient) {
	// Parse the message to get the stop id
	stopString := message.Text[1:]
	stopNumber, err := strconv.Atoi(stopString)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid stop number")
		bot.Send(msg)
		return
	}

	requestStop := apiclient.ApiApiStopsStopNumberGetRequest(client.BusAPI.ApiStopsStopNumberScheduleGet(context.Background(), int32(stopNumber)))
	stopInfo, _, err := requestStop.Execute()
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Stop not found")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Parada "+stopString+": "+*stopInfo.Name)
	bot.Send(msg)

	locationMessage := tgbotapi.NewLocation(message.Chat.ID, float64(*stopInfo.Location.Lat), float64(*stopInfo.Location.Lon))
	bot.Send(locationMessage)

	requestSchedule := apiclient.ApiApiStopsStopNumberScheduleGetRequest(client.BusAPI.ApiStopsStopNumberScheduleGet(context.Background(), int32(stopNumber)))
	schedule, _, err := requestSchedule.Execute()
	if err != nil {
		log.Println("Error getting stop schedule:", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Schedule not found")
		bot.Send(msg)
		return
	}

	text := "Horario de la parada " + stopString + ":\n\n"
	for _, scheduleItem := range schedule.Schedules {
		text += *scheduleItem.Line.Name + " | " + *scheduleItem.Route + " | " + strconv.FormatInt(int64(*scheduleItem.Time), 10) + "min\n"
	}

	for _, m := range utils.SplitLongMessage(text) {
		msg := tgbotapi.NewMessage(message.Chat.ID, m)
		bot.Send(msg)
	}
}
