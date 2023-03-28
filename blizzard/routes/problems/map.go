package problems

import (
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
)

var Map = extra.RouteMap{
	"/": {
		Methods: []models.Method{models.Get, models.Post},
		Handler: Index,
	},
	"/:problem": {
		Methods: []models.Method{models.Get, models.Patch, models.Delete},
		Handler: Problem,
	},
}
