package problems

import (
	"blizzard/blizzard/db/models/problems"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
)

func Index(ctx *extra.Context) models.Response {
	return ctx.Arr(problems.Problem{
		Id:      "hello-world",
		Title:   "Hello World",
		Contest: "",
		Tags:    []string{"easy", "beginner"},
		Content: "",
	}, problems.Problem{
		Id:      "sum",
		Title:   "Sum",
		Contest: "",
		Tags:    []string{"easy", "beginner", "numbers"},
		Content: "",
	})
}
