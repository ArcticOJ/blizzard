package oauth

import (
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
)

var Map = extra.RouteMap{
	"/": {
		Methods: []models.Method{models.Get},
		Handler: Index,
	},
	"/:provider": {
		Methods: []models.Method{models.Get},
		Handler: CreateUrl,
	},
	"/:provider/unlink": {
		Methods: []models.Method{models.Delete},
		Handler: Unlink,
	},
	"/:provider/validate": {
		Methods: []models.Method{models.Post},
		Handler: Validate,
	},
}
