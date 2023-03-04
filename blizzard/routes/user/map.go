package user

import "backend/blizzard/models"

var Map = models.RouteMap{
	"/": {
		Methods: []models.Method{models.Get, models.Post},
		Handler: Index,
	},
	"/apiKey": {
		Methods: []models.Method{models.Get, models.Patch},
		Handler: ApiKey,
	},
}
