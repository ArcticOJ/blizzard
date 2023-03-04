package feeds

import "backend/blizzard/models"

var Map = models.RouteMap{
	"/": {
		Methods: []models.Method{models.Get, models.Post},
		Handler: Index,
	},
	"/:id": {
		Methods: []models.Method{models.Get, models.Patch, models.Delete},
		Handler: func(ctx *models.Context) models.Response {
			return nil
		},
	},
}
