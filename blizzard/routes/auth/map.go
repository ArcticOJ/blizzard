package auth

import (
	"backend/blizzard/models"
	"backend/blizzard/routes/auth/oauth"
)

var Map = models.RouteMap{
	"/oauth": {
		Methods: []models.Method{models.Get},
		Handler: oauth.Index,
	},
	"/oauth/callback": {
		Methods: []models.Method{models.Get},
		Handler: oauth.Callback,
	},
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
