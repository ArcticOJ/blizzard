package users

import "blizzard/server/http"

var Map = http.RouteMap{
	"/": {
		Methods: []http.Method{http.Get, http.Delete},
		Handler: Index,
	},
}
