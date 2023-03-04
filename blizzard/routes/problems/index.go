package problems

import (
	"backend/blizzard/db/models/problems"
	"backend/blizzard/db/models/shared"
	"backend/blizzard/models"
)

func Index(ctx *models.Context) models.Response {
	author := shared.Author{
		Id:       "130139",
		Username: "AlphaNecron",
	}
	return ctx.Arr(problems.Problem{
		Id:      "hello-world",
		Title:   "Hello World",
		Contest: "",
		Tags:    []string{"easy", "beginner"},
		Author:  author,
		Content: "",
	}, problems.Problem{
		Id:      "sum",
		Title:   "Sum",
		Contest: "",
		Tags:    []string{"easy", "beginner", "numbers"},
		Author:  author,
		Content: "",
	})
}
