package submissions

import (
	"blizzard/cache/stores"
	"blizzard/models"
	"blizzard/models/extra"
	"strconv"
)

func Submission(ctx *extra.Context) models.Response {
	id := ctx.Param("submission")
	_id, e := strconv.ParseUint(id, 10, 32)
	if e != nil {
		return ctx.Bad("Invalid ID.")
	}
	if stores.Pending.IsPending(ctx.Request().Context(), uint32(_id)) {

	}
	return ctx.Respond(id)
}
