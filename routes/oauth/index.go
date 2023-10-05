package oauth

import (
	"blizzard/db"
	"blizzard/db/models/user"
	"blizzard/logger/debug"
	"blizzard/oauth"
	"blizzard/server/http"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type connection struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func Index(ctx *http.Context) http.Response {
	uid := ctx.GetUUID()
	r := echo.Map{
		"providers": oauth.EnabledProviders,
	}
	if uid != uuid.Nil {
		var c []user.OAuthConnection
		m := make(map[string]connection)
		debug.Dump(db.Database.NewSelect().Model(&c).ExcludeColumn("user_id").Where("user_id = ?", uid).Scan(ctx.Request().Context()))
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
