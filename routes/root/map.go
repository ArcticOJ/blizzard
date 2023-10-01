package root

import (
	"blizzard/models"
	"blizzard/models/extra"
)

var Map = extra.RouteMap{
	"/status": {
		Methods: []models.Method{models.Get},
		Handler: Status,
	},
}
