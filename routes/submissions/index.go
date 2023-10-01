package submissions

import (
	"blizzard/db"
	"blizzard/db/models/contest"
	"blizzard/models"
	"blizzard/models/extra"
)

func Submissions(ctx *extra.Context) models.Response {
	var submissions []contest.Submission
	if db.Database.NewSelect().Model(&submissions).Scan(ctx.Request().Context()) != nil {
		return ctx.InternalServerError("Could not fetch submissions.")
	}
	return ctx.Respond(submissions)
}
