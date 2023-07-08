package problems

import (
	"blizzard/blizzard/db/models/contest"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
)

func Problem(ctx *extra.Context) models.Response {
	return ctx.Respond(contest.Problem{
		ID:      "hello-world",
		Tags:    []string{"easy", "beginner"},
		Content: "",
	})
}
