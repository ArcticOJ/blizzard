package user

import (
	"backend/blizzard/models"
)

func Index(ctx *models.Context) models.Response {
	if user := ctx.GetUser(); user != nil {
		return ctx.Respond(user)
	}
	return ctx.NotFound("User not found")
}
