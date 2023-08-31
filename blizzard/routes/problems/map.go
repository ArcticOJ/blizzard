package problems

import (
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
)

var Map = extra.RouteMap{
	"/": {
		Methods: []models.Method{models.Get},
		Handler: Index,
	},
	"/:problem": {
		Methods: []models.Method{models.Get, models.Patch, models.Delete, models.Post},
		Handler: Problem,
	},
	"/:problem/submit": {
		Methods: []models.Method{models.Post},
		Handler: Submit,
	},
}
