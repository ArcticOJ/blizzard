package users

import (
	"github.com/ArcticOJ/blizzard/v0/cache/stores"
	"github.com/ArcticOJ/blizzard/v0/db/schema/user"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/ArcticOJ/blizzard/v0/types"
	"math"
	"strconv"
)

// GetUsers GET /
func GetUsers(ctx *http.Context) http.Response {
	page := ctx.QueryParam("page")
	var (
		p uint32 = 1
	)
	if _p, e := strconv.ParseUint(page, 10, 32); e == nil && _p > 0 {
		p = uint32(_p)
	}
	n, users := stores.Users.GetPage(ctx.Request().Context(), p)
	if n == -1 {
		return ctx.InternalServerError("Could not fetch users.")
	}
	// coerce page in 1 - max page count
	p = min(p, uint32(math.Ceil(float64(n)/float64(stores.DefaultUserPageSize))))
	//if users = stores.Users.GetPage(c, p-1, ctx.QueryParam("reversed") == "true"); users == nil {
	//	return ctx.InternalServerError("Could not fetch users.")
	//}
	return ctx.Respond(types.Paginateable[user.User]{
		Count:       uint32(n),
		CurrentPage: p,
		PageSize:    stores.DefaultUserPageSize,
		Data:        users,
	})
}
