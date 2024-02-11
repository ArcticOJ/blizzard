package user

import (
	"github.com/ArcticOJ/blizzard/v0/cache/stores"
	"github.com/ArcticOJ/blizzard/v0/server/http"
)

// GetUser GET / @auth
func GetUser(ctx *http.Context) http.Response {
	if u := stores.Users.Get(ctx.Request().Context(), ctx.GetUUID(), ""); u != nil {
		return ctx.Respond(u)
	}
	return ctx.NotFound("User not found.")
}
