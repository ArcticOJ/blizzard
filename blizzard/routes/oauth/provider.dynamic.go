package oauth

import (
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"blizzard/blizzard/oauth"
	"blizzard/blizzard/utils"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func CreateUrl(ctx *extra.Context) models.Response {
	if prov, ok := oauth.Conf[ctx.Param("provider")]; ok {
		return ctx.Respond(echo.Map{
			"url": prov.AuthCodeURL(utils.Rand(10, ""), oauth2.AccessTypeOnline),
		})
	}
	return ctx.Bad("Unsupported OAuth provider.")
}
