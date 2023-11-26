package user

import (
	"github.com/ArcticOJ/blizzard/v0/cache/stores"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/google/uuid"
	"time"
)

func HoverCard(ctx *http.Context) http.Response {
	time.Sleep(time.Second * 5)
	handle := ctx.Param("handle")
	u := stores.Users.Get(ctx.Request().Context(), uuid.Nil, handle)
	if u == nil || u.ID == uuid.Nil {
		return ctx.NotFound("User not found")
	}
	return ctx.Respond(u)
}
