package root

import (
	"blizzard/blizzard/build"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"blizzard/blizzard/server/utils"
	"github.com/labstack/echo/v4"
)

func Status(ctx *extra.Context) models.Response {
	return ctx.Respond(echo.Map{
		"version": build.Version,
		"uptime":  utils.Uptime(),
		"judges":  nil,
	})
}
