package root

import (
	"blizzard/build"
	"blizzard/cache/stores"
	"blizzard/server/http"
	"encoding/json"
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
