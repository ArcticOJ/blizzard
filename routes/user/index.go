package user

import (
	"blizzard/models"
	"blizzard/models/extra"
)

func Index(ctx *extra.Context) models.Response {
	// TODO: delete this route
	if ctx.RequireAuth() {
		return nil
	}
	if user := ctx.GetUser(); user != nil {
		return ctx.Respond(user)
	}
	return ctx.NotFound("User not found.")
}
