package users

import (
	"blizzard/db"
	"blizzard/db/models/user"
	"blizzard/models"
	"blizzard/models/extra"
	"github.com/uptrace/bun"
)

func Index(ctx *extra.Context) models.Response {
	var users []user.User
	if db.Database.NewSelect().Model(&users).Column("id", "handle", "display_name", "rating").Relation("Roles", func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Order("priority ASC").Column("icon", "color").Limit(1)
	}).Limit(50).Order("rating DESC").Scan(ctx.Request().Context()) != nil {
		return ctx.InternalServerError("Could not fetch users.")
	}
	return ctx.Respond(users)
}
