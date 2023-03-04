package auth

import "backend/blizzard/models"

func Logout(ctx *models.Context) models.Response {
	ctx.Set("user", nil)
	ctx.DeleteCookie("session")
	return ctx.Success()
}
