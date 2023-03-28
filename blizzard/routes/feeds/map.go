package feeds

import (
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
)

var Map = extra.RouteMap{
	"/": {
		Methods: []models.Method{models.Get, models.Post},
		Handler: Index,
	},
	"/:id": {
		Methods: []models.Method{models.Get, models.Patch, models.Delete},
		Handler: func(ctx *extra.Context) models.Response {
			return nil
		},
	},
}
