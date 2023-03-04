package contests

import "backend/blizzard/models"

var Map = models.RouteMap{
	"/:id/submit": {
		Methods: []models.Method{models.Post, models.Get},
		Handler: Submit,
	},
}
