package submissions

import (
	"blizzard/judge"
	"blizzard/models"
	"blizzard/models/extra"
	"strconv"
)

// TODO: broadcast cancel request to all judges
// TODO: add authorization

func CancelSubmission(ctx *extra.Context) models.Response {
	id, e := strconv.ParseUint(ctx.Param("submission"), 10, 32)
	if e != nil {
		return ctx.Bad("Invalid ID.")
	}
	return ctx.Respond(judge.ResponseObserver.Cancel(ctx.Request().Context(), uint32(id)))
}
