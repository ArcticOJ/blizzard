package submissions

import (
	"github.com/ArcticOJ/blizzard/v0/judge"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"strconv"
)

// TODO: broadcast cancel request to all judges
// TODO: add authorization

func CancelSubmission(ctx *http.Context) http.Response {
	if ctx.RequireAuth() {
		return nil
	}
	id, e := strconv.ParseUint(ctx.Param("submission"), 10, 32)
	if e != nil {
		return ctx.Bad("Invalid ID.")
	}
	return ctx.Respond(judge.Worker.Cancel(uint32(id), ctx.GetUUID()))
}
