package contests

import (
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
)

var Map = extra.RouteMap{
	"/:id/submit": {
		Methods: []models.Method{models.Post, models.Get},
		Handler: Submit,
	},
}
