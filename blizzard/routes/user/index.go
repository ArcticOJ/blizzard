package user

import (
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
)

func Index(ctx *extra.Context) models.Response {
	if ctx.RequireAuth() {
		return nil
	}
	if user := ctx.GetUser(); user != nil {
		return ctx.Respond(user)
	}
	return ctx.NotFound("User not found.")
}
