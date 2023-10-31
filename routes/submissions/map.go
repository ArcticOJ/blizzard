package submissions

import "github.com/ArcticOJ/blizzard/v0/server/http"

var Map = http.RouteMap{
	"/": {
		Methods: []http.Method{http.Get},
		Handler: Submissions,
	},
	"/:submission/source": {
		Methods: []http.Method{http.Get},
		Handler: Source,
	},
	"/:submission/cancel": {
		Methods: []http.Method{http.Get},
		Handler: CancelSubmission,
	},
	"/:submission": {
		Methods: []http.Method{http.Get},
		Handler: Submission,
	},
}
