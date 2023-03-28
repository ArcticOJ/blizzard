package auth

import (
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
)

func Logout(ctx *extra.Context) models.Response {
	ctx.Set("user", nil)
	ctx.DeleteCookie("session")
	return ctx.Success()
}
