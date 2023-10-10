package oauth

import (
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/models/user"
	"github.com/ArcticOJ/blizzard/v0/logger/debug"
	"github.com/ArcticOJ/blizzard/v0/oauth"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/ArcticOJ/blizzard/v0/utils"
)

func Unlink(ctx *http.Context) http.Response {
	if ctx.RequireAuth() {
		return nil
	}
	prov := ctx.Param("provider")
	id := ctx.QueryParam("id")
	if len(id) == 0 {
		return ctx.Bad("Invalid ID.")
	}
	if utils.ArrayIncludes(oauth.EnabledProviders, prov) {
		if _, e := db.Database.NewDelete().Model((*user.OAuthConnection)(nil)).Where("id = ? AND provider = ?", id, prov).Returning("NULL").Exec(ctx.Request().Context()); e != nil {
			debug.Dump(e)
			return ctx.NotFound("OAuth connection not found.")
		}
		return ctx.Success()
	}
	return ctx.Bad("Unsupported OAuth provider.")
}
