package users

import (
	"github.com/ArcticOJ/blizzard/v0/cache/stores"
	"github.com/ArcticOJ/blizzard/v0/db/models/user"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/ArcticOJ/blizzard/v0/types"
	"strconv"
)

func Index(ctx *http.Context) http.Response {
	page := ctx.QueryParam("page")
	var p uint16 = 1
	if _p, e := strconv.ParseUint(page, 10, 16); e == nil && _p > 0 {
		p = uint16(_p)
	}
	var users []user.MinimalUser
	if users = stores.Users.GetPage(ctx.Request().Context(), p-1, ctx.QueryParam("reversed") == "true"); users == nil {
		return ctx.InternalServerError("Could not fetch users.")
	}
	return ctx.Respond(types.Paginateable[user.MinimalUser]{
		CurrentPage: p,
		PageSize:    stores.DefaultUserPageSize,
		Data:        users,
	})
}
