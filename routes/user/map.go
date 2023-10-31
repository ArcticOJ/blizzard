package user

import "github.com/ArcticOJ/blizzard/v0/server/http"

var Map = http.RouteMap{
	"/:handle/info": {
		Methods: []http.Method{http.Get},
		Handler: Info,
	},
	"/:id/readme": {
		Methods: []http.Method{http.Get},
		Handler: Readme,
	},
	"/": {
		Methods: []http.Method{http.Get, http.Post},
		Handler: Index,
	},
	"/apiKey": {
		Methods: []http.Method{http.Get, http.Patch},
		Handler: ApiKey,
	},
	"/changeUsername": {
		Methods: []http.Method{http.Post},
		Handler: ChangeUsername,
	},
}
