package posts

import "github.com/ArcticOJ/blizzard/v0/server/http"

var Map = http.RouteMap{
	"/": {
		Methods: []http.Method{http.Get, http.Post},
		Handler: Index,
	},
	"/:id": {
		Methods: []http.Method{http.Get, http.Patch, http.Delete},
		Handler: func(ctx *http.Context) http.Response {
			return nil
		},
	},
}
