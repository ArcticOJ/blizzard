package users

import (
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/models/user"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/uptrace/bun"
)

func Index(ctx *http.Context) http.Response {
	var users []user.User
	if db.Database.NewSelect().Model(&users).Column("id", "handle", "display_name", "rating").Relation("Roles", func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Order("priority ASC").Column("icon", "color").Limit(1)
	}).Limit(50).Order("rating DESC").Scan(ctx.Request().Context()) != nil {
		return ctx.InternalServerError("Could not fetch users.")
	}
	return ctx.Respond(users)
}
