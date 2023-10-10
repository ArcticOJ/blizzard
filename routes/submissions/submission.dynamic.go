package submissions

import (
	"github.com/ArcticOJ/blizzard/v0/cache/stores"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"strconv"
)

func Submission(ctx *http.Context) http.Response {
	id := ctx.Param("submission")
	_id, e := strconv.ParseUint(id, 10, 32)
	if e != nil {
		return ctx.Bad("Invalid ID.")
	}
	if stores.Pending.IsPending(ctx.Request().Context(), uint32(_id)) {

	}
	return ctx.Respond(id)
}
