package submissions

import (
	"blizzard/db"
	"blizzard/db/models/contest"
	"blizzard/server/http"
)

func Submissions(ctx *http.Context) http.Response {
	var submissions []contest.Submission
	if db.Database.NewSelect().Model(&submissions).Scan(ctx.Request().Context()) != nil {
		return ctx.InternalServerError("Could not fetch submissions.")
	}
	return ctx.Respond(submissions)
}
