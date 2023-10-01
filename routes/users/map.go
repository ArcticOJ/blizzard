package users

import (
	"blizzard/models"
	"blizzard/models/extra"
)

var Map = extra.RouteMap{
	"/": {
		Methods: []models.Method{models.Get, models.Delete},
		Handler: Index,
	},
}
