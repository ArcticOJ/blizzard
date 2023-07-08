package user

import (
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
)

var Map = extra.RouteMap{
	// avoid conflict with other routes
	"/info/:handle": {
		Methods: []models.Method{models.Get},
		Handler: Info,
	},
	"/": {
		Methods: []models.Method{models.Get, models.Post},
		Handler: Index,
	},
	"/apiKey": {
		Methods: []models.Method{models.Get, models.Patch},
		Handler: ApiKey,
	},
}
