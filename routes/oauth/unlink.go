package oauth

import (
	"blizzard/db"
	"blizzard/db/models/user"
	"blizzard/logger/debug"
	"blizzard/models"
	"blizzard/models/extra"
	"blizzard/oauth"
	"blizzard/utils"
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
		if _, e := db.Database.NewDelete().Model((*user.OAuthConnection)(nil)).Where("id = ? AND provider = ?", id, prov).Returning("NULL").Exec(ctx.Request().Context()); e != nil {
			debug.Dump(e)
			return ctx.NotFound("OAuth connection not found.")
		}
		return ctx.Success()
	}
	return ctx.Bad("Unsupported OAuth provider.")
}
