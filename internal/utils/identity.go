package utils

import (
	apiclient "github.com/eryalito/vigo-bus-core-go-client"
)

func IsStopInFavorites(identity *apiclient.ApiIdentity, stopNumber int32) bool {
	if identity.FavoriteStops == nil {
		return false
	}
	for _, stop := range identity.FavoriteStops {
		if *stop.StopNumber == stopNumber {
			return true
		}
	}
	return false
}
