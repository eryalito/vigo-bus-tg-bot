package follow

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	apiclient "github.com/eryalito/vigo-bus-core-go-client"
	"github.com/eryalito/vigo-bus-tg-bot/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Define a struct to hold the stop, line, user, and desiredETA
type FollowTask struct {
	StopNumber int
	LineID     int
	ChatID     int64
	DesiredETA int
	Client     *apiclient.APIClient
	Bot        *tgbotapi.BotAPI
	Manager    *FollowTaskManager
}

// Function that each goroutine will run
func (task *FollowTask) Run(wg *sync.WaitGroup, interval time.Duration) {
	defer wg.Done()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// get Line info from API
	request, _, err := task.Client.BusAPI.ApiLinesGet(context.Background()).Execute()
	if err != nil {
		msg := tgbotapi.NewMessage(task.ChatID, "Error al obtener la información de la línea. Cancelado el seguimiento")
		task.Bot.Send(msg)
		task.Manager.RemoveTask(task.ChatID, task.StopNumber, task.LineID)
		return
	}

	// Check if the line exists
	var lineObject *apiclient.ApiLine
	for _, line := range request {
		if *line.Id == int32(task.LineID) {
			lineObject = &line
			break
		}
	}

	if *lineObject.Id == 0 {
		msg := tgbotapi.NewMessage(task.ChatID, "La línea no existe. Cancelado el seguimiento")
		task.Bot.Send(msg)
		task.Manager.RemoveTask(task.ChatID, task.StopNumber, task.LineID)
		return
	}

	// Check if the stop exists
	requestStop := apiclient.ApiApiStopsStopNumberGetRequest(task.Client.BusAPI.ApiStopsStopNumberScheduleGet(context.Background(), int32(task.StopNumber)))
	stopInfo, _, err := requestStop.Execute()
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(task.ChatID, "Stop not found")
		task.Bot.Send(msg)
		return
	}

	for {
		select {
		case <-ticker.C:
			currentETA, err := getCurrentETA(stopInfo, lineObject, task.Client)
			if err != nil {
				if config.Debug {
					log.Println(err)
				}
				msg := tgbotapi.NewMessage(task.ChatID, "Error al obtener el tiempo estimado de llegada. Cancelado el seguimiento")
				task.Bot.Send(msg)
				task.Manager.RemoveTask(task.ChatID, task.StopNumber, task.LineID)
				return
			}
			if config.Debug {
				log.Printf("User %d: Current ETA for stop %d, line %d is %d minutes\n", task.ChatID, task.StopNumber, task.LineID, currentETA)
			}
			if currentETA < task.DesiredETA {
				// Send message to user
				if config.Debug {
					log.Printf("User %d: Desired ETA reached for stop %d, line %d\n", task.ChatID, task.StopNumber, task.LineID)
				}

				msg := tgbotapi.NewMessage(task.ChatID, fmt.Sprintf("El bus de la línea %s ha llegado a la parada %d", *lineObject.Name, task.StopNumber))
				task.Bot.Send(msg)

				task.Manager.RemoveTask(task.ChatID, task.StopNumber, task.LineID) // Remove task from manager
				return
			}
		}
	}
}

// Dummy function to get current ETA (replace with actual implementation)
func getCurrentETA(stop *apiclient.ApiStop, line *apiclient.ApiLine, client *apiclient.APIClient) (int, error) {
	requestSchedule := apiclient.ApiApiStopsStopNumberScheduleGetRequest(client.BusAPI.ApiStopsStopNumberScheduleGet(context.Background(), *stop.StopNumber))
	schedule, _, err := requestSchedule.Execute()
	if err != nil {
		log.Println("Error getting stop schedule:", err)
		return 0, fmt.Errorf("Error al obtener horario de la línea %s en la parada %s", *line.Name, *stop.Name)
	}

	for _, scheduleItem := range schedule.Schedules {
		if *scheduleItem.Line.Id == *line.Id {
			return int(*scheduleItem.Time), nil
		}
	}

	return 0, fmt.Errorf("No se ha encontrado la línea %s en el horario de la parada %s", *line.Name, *stop.Name)
}
