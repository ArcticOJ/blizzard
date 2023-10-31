package submissions

import (
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/models/contest"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/ArcticOJ/blizzard/v0/storage"
)

func Source(ctx *http.Context) http.Response {
	if ctx.RequireAuth() {
		return nil
	}
	id := ctx.Param("submission")
	s := new(contest.Submission)
	if db.Database.NewSelect().Model(s).Where("id = ?", id).Column("id", "author_id", "extension", "problem_id").Scan(ctx.Request().Context()) != nil {
		return ctx.NotFound("Submission not found.")
	}
	if s.AuthorID != ctx.GetUUID() {
		return ctx.Unauthorized()
	}
	if e := ctx.Inline(storage.Submission.GetPath(s.ID, s.Extension), s.ProblemID+"."+s.Extension); e != nil {
		return ctx.InternalServerError("Failed to load source code.")
	}
	return nil
}
