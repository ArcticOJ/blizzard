package submissions

import (
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/schema/contest"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/ArcticOJ/blizzard/v0/storage"
	"path"
)

func Source(ctx *http.Context) http.Response {
	if ctx.RequireAuth() {
		return nil
	}
	id := ctx.Param("submission")
	s := new(contest.Submission)
	if db.Database.NewSelect().Model(s).Where("id = ?", id).Column("id", "author_id", "file_name", "problem_id").Scan(ctx.Request().Context()) != nil {
		return ctx.NotFound("Submission not found.")
	}
	if s.AuthorID != ctx.GetUUID() {
		return ctx.Unauthorized()
	}
	downloadFileName := s.ProblemID + "." + path.Ext(s.FileName)
	if e := ctx.Inline(storage.Submission.GetPath(s.FileName), downloadFileName); e != nil {
		return ctx.InternalServerError("Failed to load source code.")
	}
	return nil
}
