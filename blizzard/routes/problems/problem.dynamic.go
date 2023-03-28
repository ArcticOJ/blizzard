package problems

import (
	"blizzard/blizzard/db/models/problems"
	"blizzard/blizzard/db/models/shared"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
)

func Problem(ctx *extra.Context) models.Response {
	return ctx.Respond(problems.Problem{
		Id:      "hello-world",
		Title:   "Hello World",
		Contest: "",
		Tags:    []string{"easy", "beginner"},
		Author: shared.Author{
			Id:       "130139",
			Username: "AlphaNecron",
		},
		Content: "",
	})
}
