package auth

import "github.com/ArcticOJ/blizzard/v0/server/http"

var Map = http.RouteMap{
	"/login": {
		Methods: []http.Method{http.Post},
		Handler: Login,
	},
	"/register": {
		Methods: []http.Method{http.Post},
		Handler: Register,
	},
	"/logout": {
		Methods: []http.Method{http.Get},
		Handler: Logout,
	},
}
