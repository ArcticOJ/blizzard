package user

import (
	"blizzard/cache/stores"
	"blizzard/models"
	"blizzard/models/extra"
	"github.com/google/uuid"
)

// TODO: allow querying by uuid

func Info(ctx *extra.Context) models.Response {
	handle := ctx.Param("handle")
	u := stores.Users.Get(ctx.Request().Context(), uuid.Nil, handle)
	if u == nil || u.ID == uuid.Nil {
		return ctx.NotFound("User not found")
	}
	return ctx.Respond(u)
}
