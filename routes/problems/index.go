package problems

import (
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/schema/contest"
	"github.com/ArcticOJ/blizzard/v0/server/http"
)

// TODO: add filters and pagination

func Index(ctx *http.Context) http.Response {
	var problems []contest.Problem
	if db.Database.NewSelect().Model(&problems).Column("id", "tags", "title", "author_id").Scan(ctx.Request().Context()) != nil {
		return ctx.InternalServerError("Could not fetch problems.")
	}
	return ctx.Respond(problems)
}
