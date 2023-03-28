package oauth

import (
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"blizzard/blizzard/oauth"
)

func Index(ctx *extra.Context) models.Response {
	return ctx.Respond(oauth.EnabledProviders)
}
