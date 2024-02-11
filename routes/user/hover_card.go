package user

import (
	"github.com/ArcticOJ/blizzard/v0/cache/stores"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/google/uuid"
)

// GetHoverCard GET /:handle/hoverCard
func GetHoverCard(ctx *http.Context) http.Response {
	handle := ctx.Param("handle")
	u := stores.Users.Get(ctx.Request().Context(), uuid.Nil, handle)
	if u == nil || u.ID == uuid.Nil {
		return ctx.NotFound("User not found")
	}
	return ctx.Respond(u)
}
