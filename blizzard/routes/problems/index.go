package problems

import (
	"blizzard/blizzard/db"
	"blizzard/blizzard/db/models/contest"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
)

// TODO: add filters and pagination

func Index(ctx *extra.Context) models.Response {
	var problems []contest.Problem
	if db.Database.NewSelect().Model(&problems).Column("id", "tags", "title", "author_id").Scan(ctx.Request().Context()) != nil {
		return ctx.InternalServerError("Could not fetch problems.")
	}
	return ctx.Respond(problems)
}
