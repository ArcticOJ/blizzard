package user

import (
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/ArcticOJ/blizzard/v0/storage"
	"github.com/google/uuid"
)

// GetReadme GET /:id/readme
func GetReadme(ctx *http.Context) http.Response {
	uid, e := uuid.Parse(ctx.Param("id"))
	if e != nil {
		return ctx.Bad("Failed to parse UUID.")
	}
	if e = ctx.Inline(storage.READMEs.GetPath(uid), "readme.md"); e != nil {
		return ctx.InternalServerError("Failed to load README.")
	}
	return nil
}
