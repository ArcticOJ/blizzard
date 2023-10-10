package root

import "github.com/ArcticOJ/blizzard/v0/server/http"

var Map = http.RouteMap{
	"/status": {
		Methods: []http.Method{http.Get},
		Handler: Status,
	},
}
