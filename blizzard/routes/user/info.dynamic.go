package user

import (
	"blizzard/blizzard/cache"
	"blizzard/blizzard/db"
	"blizzard/blizzard/db/models/user"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"github.com/uptrace/bun"
	"time"
)

// TODO: allow querying by uuid

func Info(ctx *extra.Context) models.Response {
	handle := ctx.Param("handle")
	usr := cache.Get[user.User](ctx.Request().URL.Path, func() (u user.User, e error) {
		e = db.Database.NewSelect().Model(&u).Where("handle = ?", handle).Relation("Connections", func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.Where("show_in_profile = true")
		}).Relation("Roles", func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.Order("priority ASC").Column("name", "icon", "style", "name_style")
		}).ExcludeColumn("email_verified", "api_key", "password").Scan(ctx.Request().Context())
		return
	}, time.Minute*5)
	if usr == nil {
		return ctx.NotFound("User not found.")
	}
	return ctx.Respond(usr)
}
