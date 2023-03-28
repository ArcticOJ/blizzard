package problems

import (
	"blizzard/blizzard/db/models/problems"
	"blizzard/blizzard/db/models/shared"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
)

func Index(ctx *extra.Context) models.Response {
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
