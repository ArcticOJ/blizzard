package oauth

import (
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/schema/user"
	"github.com/ArcticOJ/blizzard/v0/logger/debug"
	"github.com/ArcticOJ/blizzard/v0/server/http"
)

type connection struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// GetConnections GET /connections @auth
func GetConnections(ctx *http.Context) http.Response {
	var c []user.OAuthConnection
	m := make(map[string]connection)
	debug.Dump(db.Database.NewSelect().Model(&c).ExcludeColumn("user_id").Where("user_id = ?", ctx.GetUUID()).Scan(ctx.Request().Context()))
	debug.Dump(c)
	for _, p := range c {
		m[p.Provider] = connection{
			Username: p.Username,
		}
	}
	return ctx.Respond(m)
}
