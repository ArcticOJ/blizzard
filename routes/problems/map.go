package problems

import "blizzard/server/http"

var Map = http.RouteMap{
	"/": {
		Methods: []http.Method{http.Get},
		Handler: Index,
	},
	"/:problem": {
		Methods: []http.Method{http.Get, http.Patch, http.Delete, http.Post},
		Handler: Problem,
	},
	"/:problem/submit": {
		Methods: []http.Method{http.Post},
		Handler: Submit,
	},
}
