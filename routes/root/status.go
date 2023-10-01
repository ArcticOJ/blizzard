package root

import (
	"blizzard/build"
	"blizzard/judge"
	"blizzard/models"
	"blizzard/models/extra"
	"github.com/labstack/echo/v4"
)

func Status(ctx *extra.Context) models.Response {
	return ctx.Respond(echo.Map{
		"version": build.Version,
		"judges":  judge.GetStatus(),
	})
}
