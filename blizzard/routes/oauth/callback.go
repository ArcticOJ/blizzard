package oauth

import (
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"blizzard/blizzard/oauth"
	"fmt"
)

func Callback(ctx *extra.Context) models.Response {
	// TODO: handle oauth callback
	if prov, ok := oauth.Conf[ctx.Param("provider")]; ok {
		code, state := ctx.QueryParam("code"), ctx.QueryParam("state")
		if code == "" || state == "" {
			return ctx.Bad("Invalid code or state.")
		}
		fmt.Println(prov.Exchange(ctx.Request().Context(), code))
		return ctx.Success()
	}
	return ctx.Bad("Unsupported OAuth provider.")
}
