package submissions

import (
	"blizzard/judge"
	"blizzard/server/http"
	"strconv"
)

// TODO: broadcast cancel request to all judges
// TODO: add authorization

func CancelSubmission(ctx *http.Context) http.Response {
	id, e := strconv.ParseUint(ctx.Param("submission"), 10, 32)
	if e != nil {
		return ctx.Bad("Invalid ID.")
	}
	return ctx.Respond(judge.ResponseWorker.Cancel(ctx.Request().Context(), uint32(id)))
}
