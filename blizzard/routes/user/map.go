package user

import (
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
)

var Map = extra.RouteMap{
	"/": {
		Methods: []models.Method{models.Get, models.Post},
		Handler: Index,
	},
	"/apiKey": {
		Methods: []models.Method{models.Get, models.Patch},
		Handler: ApiKey,
	},
}
