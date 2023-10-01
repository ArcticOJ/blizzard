package problems

import (
	"blizzard/db"
	"blizzard/db/models/contest"
	"blizzard/models"
	"blizzard/models/extra"
)

func Problem(ctx *extra.Context) models.Response {
	id := ctx.Param("problem")
	var prob contest.Problem
	if db.Database.NewSelect().Model(&prob).Where("problem.id = ?", id).Scan(ctx.Request().Context()) != nil {
		return ctx.NotFound("Problem not found.")
	}
	return ctx.Respond(prob)
}
