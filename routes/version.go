package routes

import (
	"github.com/ArcticOJ/blizzard/v0/build"
	"github.com/ArcticOJ/blizzard/v0/config"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/labstack/echo/v4"
)

// GetVersion GET /version
func GetVersion(ctx *http.Context) http.Response {
	return ctx.Respond(echo.Map{
		"brand":     config.Config.Brand,
		"buildDate": build.Date,
		"buildHash": build.Hash,
		"version":   build.Version,
	})
}
