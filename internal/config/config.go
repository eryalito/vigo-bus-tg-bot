package config

import (
	"flag"
	"os"
	"strconv"
)

var (
	APIURL      string
	APIProtocol string
	APIToken    string
	BotToken    string
	Debug       bool
)

func Init() {
	// Define command-line flags
	flag.StringVar(&APIURL, "api-url", getEnv("API_URL", "localhost:8080"), "URL of the API server")
	flag.StringVar(&APIProtocol, "api-protocol", getEnv("API_PROTOCOL", "http"), "Protocol to use")
	flag.StringVar(&APIToken, "api-token", getEnv("API_TOKEN", "your-secret-token"), "Authentication token")
	flag.StringVar(&BotToken, "bot-token", getEnv("BOT_TOKEN", "bot-token"), "Telegram bot token")

	// Convert the DEBUG environment variable to a bool
	debugEnv := getEnv("DEBUG", "false")
	debugValue, err := strconv.ParseBool(debugEnv)
	if err != nil {
		debugValue = false
	}
	flag.BoolVar(&Debug, "debug", debugValue, "Enable debug mode")

	// Parse command-line flags
	flag.Parse()
}

// getEnv reads an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
