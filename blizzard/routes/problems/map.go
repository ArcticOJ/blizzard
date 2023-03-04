package problems

import "backend/blizzard/models"

var Map = models.RouteMap{
	"/": {
		Methods: []models.Method{models.Get, models.Post},
		Handler: Index,
	},
	"/:problem": {
		Methods: []models.Method{models.Get, models.Patch, models.Delete},
		Handler: Problem,
	},
}
