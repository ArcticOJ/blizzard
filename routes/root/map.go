package root

import "blizzard/server/http"

var Map = http.RouteMap{
	"/status": {
		Methods: []http.Method{http.Get},
		Handler: Status,
	},
}
