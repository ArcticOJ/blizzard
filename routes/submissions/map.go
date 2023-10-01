package submissions

import (
	"blizzard/models"
	"blizzard/models/extra"
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
