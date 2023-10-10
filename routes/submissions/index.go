package submissions

import (
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/models/contest"
	"github.com/ArcticOJ/blizzard/v0/server/http"
)

func Submissions(ctx *http.Context) http.Response {
	var submissions []contest.Submission
	if db.Database.NewSelect().Model(&submissions).Scan(ctx.Request().Context()) != nil {
		return ctx.InternalServerError("Could not fetch submissions.")
	}
	return ctx.Respond(submissions)
}
