package user

import (
	"blizzard/cache/stores"
	"blizzard/server/http"
)

func Index(ctx *http.Context) http.Response {
	if ctx.RequireAuth() {
		return nil
	}
	if user := stores.Users.GetMinimal(ctx.Request().Context(), ctx.GetUUID()); user != nil {
		return ctx.Respond(user)
	}
	return ctx.NotFound("User not found.")
}
