package root

import (
	"blizzard/build"
	"blizzard/judge"
	"blizzard/server/http"
	"github.com/labstack/echo/v4"
)

func Status(ctx *http.Context) http.Response {
	return ctx.Respond(echo.Map{
		"version": build.Version,
		"judges":  judge.GetStatus(),
	})
}
