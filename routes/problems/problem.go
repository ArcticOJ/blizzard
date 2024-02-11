package problems

import (
	"github.com/ArcticOJ/blizzard/v0/cache/stores"
	"github.com/ArcticOJ/blizzard/v0/server/http"
)

// GetProblem GET /:id
func GetProblem(ctx *http.Context) http.Response {
	id := ctx.Param("problem")
	if p := stores.Problems.Get(ctx.Request().Context(), id); p != nil {
		return ctx.Respond(p)
	}
	return ctx.NotFound("Problem not found.")
}
