package user

import "blizzard/server/http"

func Index(ctx *http.Context) http.Response {
	// TODO: delete this route
	if ctx.RequireAuth() {
		return nil
	}
	if user := ctx.GetUser(); user != nil {
		return ctx.Respond(user)
	}
	return ctx.NotFound("User not found.")
}
