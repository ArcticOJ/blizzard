package auth

import (
	"blizzard/models"
	"blizzard/models/extra"
)

func Logout(ctx *extra.Context) models.Response {
	ctx.Set("user", nil)
	ctx.DeleteCookie("session")
	return ctx.Success()
}
