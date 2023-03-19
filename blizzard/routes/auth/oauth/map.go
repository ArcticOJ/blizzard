package oauth

import "backend/blizzard/models"

var Map = models.RouteMap{
	"/oauth/list": {
		Methods: []models.Method{models.Get},
		Handler: List,
	},
	"/oauth/:provider": {
		Methods: []models.Method{models.Get},
		Handler: Index,
	},
	"/oauth/callback/:provider": {
		Methods: []models.Method{models.Get},
		Handler: Callback,
	},
}
