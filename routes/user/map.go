package user

import "blizzard/server/http"

var Map = http.RouteMap{
	// avoid conflict with other routes
	"/info/:handle": {
		Methods: []http.Method{http.Get},
		Handler: Info,
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
