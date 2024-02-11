package submissions

import (
	"github.com/ArcticOJ/blizzard/v0/judge"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"strconv"
)

// TODO: broadcast cancel request to all judges
// TODO: add authorization

// CancelSubmission POST /:id/cancel @auth
func CancelSubmission(ctx *http.Context) http.Response {
	id, e := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if e != nil {
		return ctx.Bad("Invalid ID.")
	}
	return ctx.Respond(judge.Observer.Cancel(uint32(id), ctx.GetUUID()))
}
