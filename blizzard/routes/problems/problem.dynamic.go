package problems

import (
	"backend/blizzard/db/models/problems"
	"backend/blizzard/db/models/shared"
	"backend/blizzard/models"
)

func Problem(ctx *models.Context) models.Response {
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
