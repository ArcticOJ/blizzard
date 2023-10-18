package auth

import (
	"github.com/ArcticOJ/blizzard/v0/server/http"
)

func Logout(ctx *http.Context) http.Response {
	ctx.Set("id", nil)
	ctx.DeleteCookie("session")
	return ctx.Success()
}
