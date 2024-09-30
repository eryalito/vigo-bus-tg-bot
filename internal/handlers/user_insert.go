package handlers

import (
	"context"
	"log"
	"strconv"

	apiclient "github.com/eryalito/vigo-bus-core-go-client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func InsertHandler(identity *apiclient.ApiIdentity, bot *tgbotapi.BotAPI, chatID int64, client *apiclient.APIClient) (*apiclient.ApiIdentity, error) {

	requestGet := client.IdentityAPI.ApiUsersProviderUuidGet(context.Background(), "telegram", strconv.FormatInt(chatID, 10))
	responseIdentity, _, err := requestGet.Execute()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if responseIdentity != nil {
		// User already exists, do nothing
		return responseIdentity, nil
	}

	request := client.IdentityAPI.ApiUsersProviderUuidPost(context.Background(), "telegram", strconv.FormatInt(chatID, 10))

	responseIdentity, _, err = request.Execute()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return responseIdentity, nil

}
