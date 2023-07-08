package problems

import (
	"blizzard/blizzard/db/models/contest"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
)

func Index(ctx *extra.Context) models.Response {
	return ctx.Arr(contest.Problem{
		ID:      "hello-world",
		Tags:    []string{"easy", "beginner"},
		Content: "",
	}, contest.Problem{
		ID:      "sum",
		Tags:    []string{"easy", "beginner", "numbers"},
		Content: "",
	})
}
