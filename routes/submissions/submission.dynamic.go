package submissions

import (
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/schema/contest"
	"github.com/ArcticOJ/blizzard/v0/judge"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/ArcticOJ/polar/v0/types"
	"github.com/uptrace/bun"
)

func Submission(ctx *http.Context) http.Response {
	if ctx.RequireAuth() {
		return nil
	}
	id := ctx.Param("submission")
	s := new(contest.Submission)
	if db.Database.NewSelect().Model(s).Where("submission.id = ?", id).Relation("Problem", func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Column("id", "title")
	}).Relation("Author", func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Column("handle", "id")
	}).Scan(ctx.Request().Context()) != nil {
		return ctx.NotFound("Submission not found.")
	}
	if s.AuthorID != ctx.GetUUID() {
		return ctx.Unauthorized()
	}
	resChan, element := judge.Observer.Subscribe(types.Submission{
		ID:            s.ID,
		TestCount:     s.Problem.TestCount,
		PointsPerTest: s.Problem.PointsPerTest,
	}, func() interface{} {
		return judge.Observer.GetResults(s.ID)
	})
	if resChan == nil {
		return ctx.Respond(s)
	}
	stream := ctx.StreamResponse()
	done := ctx.Request().Context().Done()
	for {
		select {
		case <-done:
			judge.Observer.Unsubscribe(s.ID, element)
			return nil
		case r, more := <-resChan:
			if !more || stream.Write(r) != nil {
				return nil
			}
		}
	}
}
