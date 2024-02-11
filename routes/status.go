package routes

import (
	"github.com/ArcticOJ/blizzard/v0/judge"
	"github.com/ArcticOJ/blizzard/v0/server/http"
)

// GetStatus GET /status
func GetStatus(ctx *http.Context) http.Response {
	return ctx.Respond(judge.Observer.GetJudges())
}
