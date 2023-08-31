package submissions

import (
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
)

var Map = extra.RouteMap{
	"/": {
		Methods: []models.Method{models.Get},
		Handler: Submissions,
	},
	"/:submission/cancel": {
		Methods: []models.Method{models.Get},
		Handler: CancelSubmission,
	},
	"/:submission": {
		Methods: []models.Method{models.Get},
		Handler: Submission,
	},
}
