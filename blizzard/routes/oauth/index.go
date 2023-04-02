package oauth

import (
	"blizzard/blizzard/db"
	"blizzard/blizzard/db/models/users"
	"blizzard/blizzard/logger/debug"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"blizzard/blizzard/oauth"
	"github.com/labstack/echo/v4"
)

type connection struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

func Index(ctx *extra.Context) models.Response {
	uuid := ctx.GetUUID()
	r := echo.Map{
		"providers": oauth.EnabledProviders,
	}
	if uuid != nil {
		var c []users.OAuthConnection
		m := make(map[string]connection)
		debug.Dump(db.Database.NewSelect().Model((*users.OAuthConnection)(nil)).ExcludeColumn("user_id").Where("user_id = ?", uuid).Scan(ctx.Request().Context(), &c))
		debug.Dump(c)
		for _, p := range c {
			m[p.Provider] = connection{
				Id:       p.ID,
				Username: p.Username,
			}
		}
		r["connections"] = m
	}
	return ctx.Respond(r)
}
