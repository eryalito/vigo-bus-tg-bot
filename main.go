package main

import (
	"log"
	"regexp"
	"strings"

	apiclient "github.com/eryalito/vigo-bus-core-go-client"
	"github.com/eryalito/vigo-bus-tg-bot/internal/api"
	"github.com/eryalito/vigo-bus-tg-bot/internal/config"
	"github.com/eryalito/vigo-bus-tg-bot/internal/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	log.Println("Starting Vigo Bus Telegram Bot")
	log.Println("Initializing configuration")
	config.Init()
	log.Println("Configuration initialized")

	client := api.GetAPIClient()

	// Initialize the bot
	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	bot.Debug = config.Debug
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		go processUpdate(update, bot, client)
	}

}

func processUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI, client *apiclient.APIClient) {
	if update.Message != nil { // If we got a message
		if config.Debug {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		}
		identity, err := handlers.InsertHandler(nil, bot, update.Message, client)
		if err != nil || identity == nil {
			log.Println(err)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to process the message"))
			return
		}
		// Handle commands
		switch {
		case strings.HasPrefix(update.Message.Text, "/help"):
			handlers.HelpHandler(identity, bot, update.Message, client)
		case strings.HasPrefix(update.Message.Text, "/start"):
			handlers.StartHandler(identity, bot, update.Message, client)
		case strings.HasPrefix(update.Message.Text, "/search"):
			handlers.FindHandler(identity, bot, update.Message, client)
		case strings.HasPrefix(update.Message.Text, "/fav"):
			handlers.FavoriteHandler(identity, bot, update.Message, client)
		case matchRegex(update.Message.Text, `^/\d+$`):
			handlers.StopInfoHandler(identity, bot, update.Message, client)
		default:
			if update.Message.Location != nil {
				handlers.LocationHandler(identity, bot, update.Message, client)
				return
			}
			handlers.UnknownHandler(identity, bot, update.Message, client)
		}
	}
}

// matchRegex checks if the text matches the given regex pattern
func matchRegex(text, pattern string) bool {
	matched, err := regexp.MatchString(pattern, text)
	if err != nil {
		return false
	}
	return matched
}
