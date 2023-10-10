package problems

import (
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/models/contest"
	"github.com/ArcticOJ/blizzard/v0/server/http"
)

func Problem(ctx *http.Context) http.Response {
	id := ctx.Param("problem")
	var prob contest.Problem
	if db.Database.NewSelect().Model(&prob).Where("problem.id = ?", id).Scan(ctx.Request().Context()) != nil {
		return ctx.NotFound("Problem not found.")
	}
	return ctx.Respond(prob)
}
