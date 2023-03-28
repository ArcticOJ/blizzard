package root

import (
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
)

var Map = extra.RouteMap{
	"/health": {
		Methods: []models.Method{models.Get},
		Handler: Health,
	},
}
