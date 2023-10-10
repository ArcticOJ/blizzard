package root

import (
	"encoding/json"
	"github.com/ArcticOJ/blizzard/v0/build"
	"github.com/ArcticOJ/blizzard/v0/cache/stores"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/labstack/echo/v4"
)

func Status(ctx *http.Context) http.Response {
	c := ctx.Request().Context()
	return ctx.Respond(echo.Map{
		"version": build.Version,
		"status":  json.RawMessage(stores.Judge.GetJudgeStatus(c)),
		"judges":  stores.Judge.GetJudgeList(c),
	})
}
