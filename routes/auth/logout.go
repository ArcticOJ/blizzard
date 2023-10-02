package auth

import (
	"blizzard/server/http"
)

func Logout(ctx *http.Context) http.Response {
	ctx.Set("user", nil)
	ctx.DeleteCookie("session")
	return ctx.Success()
}
