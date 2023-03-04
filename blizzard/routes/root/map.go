package root

import "backend/blizzard/models"

var Map = models.RouteMap{
	"/health": {
		Methods: []models.Method{models.Get},
		Handler: Health,
	},
}
