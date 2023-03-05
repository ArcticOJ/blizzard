package oauth

import "backend/blizzard/models"

var Map = models.RouteMap{
	"/oauth/:provider": {
		Methods: []models.Method{models.Get},
		Handler: Index,
	},
	"/oauth/callback/:provider": {
		Methods: []models.Method{models.Get},
		Handler: Callback,
	},
}
