package users

import "github.com/ArcticOJ/blizzard/v0/server/http"

var Map = http.RouteMap{
	"/": {
		Methods: []http.Method{http.Get, http.Delete},
		Handler: Index,
	},
}
