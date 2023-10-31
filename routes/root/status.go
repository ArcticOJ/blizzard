package root

import (
	"github.com/ArcticOJ/blizzard/v0/judge"
	"github.com/ArcticOJ/blizzard/v0/server/http"
)

func Status(ctx *http.Context) http.Response {
	return ctx.Respond(judge.Worker.GetStatus())
}
