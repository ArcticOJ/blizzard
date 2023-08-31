package submissions

import (
	"blizzard/blizzard/cache"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"strconv"
)

func Submission(ctx *extra.Context) models.Response {
	id := ctx.Param("submission")
	_id, e := strconv.ParseUint(id, 10, 32)
	if e != nil {
		return ctx.Bad("Invalid ID.")
	}
	if cache.Result.IsPending(uint32(_id)) {

	}
	return ctx.Respond(id)
}
