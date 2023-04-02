package oauth

import (
	"blizzard/blizzard/db"
	"blizzard/blizzard/db/models/users"
	"blizzard/blizzard/logger/debug"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"blizzard/blizzard/oauth"
	"blizzard/blizzard/utils"
)

func Unlink(ctx *extra.Context) models.Response {
	if ctx.RequireAuth() {
		return nil
	}
	prov := ctx.Param("provider")
	id := ctx.QueryParam("id")
	if len(id) == 0 {
		return ctx.Bad("Invalid ID.")
	}
	if utils.ArrayIncludes(oauth.EnabledProviders, prov) {
		if _, e := db.Database.NewDelete().Model((*users.OAuthConnection)(nil)).Where("id = ? AND provider = ?", id, prov).Returning("NULL").Exec(ctx.Request().Context()); e != nil {
			debug.Dump(e)
			return ctx.NotFound("OAuth connection not found.")
		}
		return ctx.Success()
	}
	return ctx.Bad("Unsupported OAuth provider.")
}
