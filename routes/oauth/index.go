package oauth

import (
	"blizzard/db"
	"blizzard/db/models/user"
	"blizzard/logger/debug"
	"blizzard/models"
	"blizzard/models/extra"
	"blizzard/oauth"
	"github.com/labstack/echo/v4"
)

type connection struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func Index(ctx *extra.Context) models.Response {
	uuid := ctx.GetUUID()
	r := echo.Map{
		"providers": oauth.EnabledProviders,
	}
	if uuid != nil {
		var c []user.OAuthConnection
		m := make(map[string]connection)
		debug.Dump(db.Database.NewSelect().Model(&c).ExcludeColumn("user_id").Where("user_id = ?", uuid).Scan(ctx.Request().Context()))
		debug.Dump(c)
		for _, p := range c {
			m[p.Provider] = connection{
				ID:       p.ID,
				Username: p.Username,
			}
		}
		r["connections"] = m
	}
	return ctx.Respond(r)
}
