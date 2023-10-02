package user

import (
	"blizzard/cache/stores"
	"blizzard/server/http"
	"github.com/google/uuid"
)

// TODO: allow querying by uuid

func Info(ctx *http.Context) http.Response {
	handle := ctx.Param("handle")
	u := stores.Users.Get(ctx.Request().Context(), uuid.Nil, handle)
	if u == nil || u.ID == uuid.Nil {
		return ctx.NotFound("User not found")
	}
	return ctx.Respond(u)
}
