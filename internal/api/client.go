package api

import (
	"sync"

	apiclient "github.com/eryalito/vigo-bus-core-go-client"
	"github.com/eryalito/vigo-bus-tg-bot/internal/config"
)

var (
	clientInstance *apiclient.APIClient
	once           sync.Once
)

// GetAPIClient returns the singleton instance of the API client
func GetAPIClient() *apiclient.APIClient {
	once.Do(func() {
		cfg := apiclient.NewConfiguration()
		cfg.Host = config.APIURL // Replace with the actual host
		cfg.Servers = apiclient.ServerConfigurations{
			{
				URL: "http://" + config.APIURL,
			},
		} // Replace with the actual host
		cfg.AddDefaultHeader("Authorization", "Bearer "+config.APIToken) // Replace with your actual API key
		clientInstance = apiclient.NewAPIClient(cfg)
	})
	return clientInstance
}
