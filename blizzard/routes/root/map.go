package root

import (
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
)

var Map = extra.RouteMap{
	"/status": {
		Methods: []models.Method{models.Get},
		Handler: Status,
	},
}
