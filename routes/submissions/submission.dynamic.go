package submissions

import (
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/models/contest"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/uptrace/bun"
)

func Submission(ctx *http.Context) http.Response {
	id := ctx.Param("submission")
	s := new(contest.Submission)
	if db.Database.NewSelect().Model(s).Where("submission.id = ?", id).Relation("Problem", func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Column("id", "title")
	}).Relation("Author", func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Column("handle", "id")
	}).Scan(ctx.Request().Context()) != nil {
		return ctx.NotFound("Submission not found.")
	}
	return ctx.Respond(s)
	//_id, e := strconv.ParseUint(id, 10, 32)
	//if e != nil {
	//	return ctx.Bad("Invalid ID.")
	//}
	//if stores.Submissions.IsPending(ctx.Request().Context(), uint32(_id)) {
	//
	//}
	//return ctx.Respond(id)
}
