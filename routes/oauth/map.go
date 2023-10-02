package oauth

import "blizzard/server/http"

var Map = http.RouteMap{
	"/": {
		Methods: []http.Method{http.Get},
		Handler: Index,
	},
	"/:provider": {
		Methods: []http.Method{http.Get},
		Handler: CreateUrl,
	},
	"/:provider/unlink": {
		Methods: []http.Method{http.Delete},
		Handler: Unlink,
	},
	"/:provider/validate": {
		Methods: []http.Method{http.Post},
		Handler: Validate,
	},
}
