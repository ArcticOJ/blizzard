package problems

import (
	"blizzard/blizzard/cache"
	"blizzard/blizzard/db"
	"blizzard/blizzard/db/models/contest"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"time"
)

func Problem(ctx *extra.Context) models.Response {
	id := ctx.Param("problem")
	problem := cache.Get[contest.Problem](ctx.Request().URL.Path, func() (prob contest.Problem, e error) {
		e = db.Database.NewSelect().Model(&prob).Where("problem.id = ?", id).Scan(ctx.Request().Context())
		return
	}, time.Minute*5)
	if problem == nil {
		return ctx.NotFound("Problem not found.")
	}
	return ctx.Respond(problem)
}
