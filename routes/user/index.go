package user

import (
	"github.com/ArcticOJ/blizzard/v0/server/http"
)

func Index(ctx *http.Context) http.Response {
	if ctx.RequireAuth() {
		return nil
	}
	if u := ctx.GetUser(); u != nil {
		return ctx.Respond(u)
	}
	return ctx.NotFound("User not found.")
}
