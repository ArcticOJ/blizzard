package auth

import (
	"backend/blizzard/models"
)

var Map = models.RouteMap{
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
