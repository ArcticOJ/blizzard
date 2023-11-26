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
	var (
		p     uint32 = 1
		users []user.MinimalUser
	)
	if _p, e := strconv.ParseUint(page, 10, 32); e == nil && _p > 0 {
		p = uint32(_p)
	}
	c := ctx.Request().Context()
	if users = stores.Users.GetPage(c, p-1, ctx.QueryParam("reversed") == "true"); users == nil {
		return ctx.InternalServerError("Could not fetch users.")
	}
	return ctx.Respond(types.Paginateable[user.MinimalUser]{
		Count:       stores.Users.UserCount(c),
		CurrentPage: p,
		PageSize:    stores.DefaultUserPageSize,
		Data:        users,
	})
}
