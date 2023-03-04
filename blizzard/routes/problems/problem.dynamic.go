package problems

import (
	"backend/blizzard/db/models/problems"
	"backend/blizzard/db/models/shared"
	"backend/blizzard/models"
	"fmt"
)

func Problem(ctx *models.Context) models.Response {
	fmt.Println(ctx.Param("problem"))
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
