package auth

import (
	"blizzard/models"
	"blizzard/models/extra"
)

var Map = extra.RouteMap{
	"/login": {
		Methods: []models.Method{models.Post},
		Handler: Login,
	},
	"/register": {
		Methods: []models.Method{models.Post},
		Handler: Register,
	},
	"/logout": {
		Methods: []models.Method{models.Get},
		Handler: Logout,
	},
}
