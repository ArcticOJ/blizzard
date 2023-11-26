// TODO: complete submit route

package problems

import (
	"container/list"
	"context"
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/models/contest"
	"github.com/ArcticOJ/blizzard/v0/judge"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/ArcticOJ/blizzard/v0/storage"
	"github.com/ArcticOJ/polar/v0/types"
	"github.com/google/uuid"
	"path"
	"slices"
	"strings"
)

func prepare(id uint32, _path, runtime string, userId uuid.UUID, problem *contest.Problem) types.Submission {
	c := problem.Constraints
	return types.Submission{
		AuthorID:      userId.String(),
		ID:            id,
		Runtime:       runtime,
		SourcePath:    path.Base(_path),
		ProblemID:     problem.ID,
		TestCount:     problem.TestCount,
		PointsPerTest: problem.PointsPerTest,
		Constraints: types.Constraints{
			IsInteractive: c.IsInteractive,
			TimeLimit:     c.TimeLimit,
			MemoryLimit:   c.MemoryLimit,
			OutputLimit:   c.OutputLimit,
			AllowPartial:  c.AllowPartial,
			ShortCircuit:  c.ShortCircuit,
		},
	}
}

func getExt(fileName string) string {
	return strings.ToLower(strings.TrimLeft(path.Ext(fileName), "."))
}

func createSubmission(ctx context.Context, userId uuid.UUID, problem, runtime string, ext string) (*contest.Submission, func() error, func() error) {
	sub := &contest.Submission{
		AuthorID:  userId,
		ProblemID: problem,
		Runtime:   runtime,
		Extension: ext,
	}
	tx, e := db.Database.Begin()
	if e != nil {
		return nil, nil, nil
	}
	if _, e = tx.NewInsert().Model(sub).Returning("id, submitted_at", "author_id").Exec(ctx); e != nil {
		tx.Rollback()
		return nil, nil, nil
	}
	return sub, tx.Rollback, tx.Commit
}

func Submit(ctx *http.Context) http.Response {
	var (
		res     chan interface{}
		element *list.Element
		e       error
	)
	if ctx.RequireAuth() {
		return nil
	}
	code, e := ctx.FormFile("code")
	shouldStream := ctx.FormValue("stream") == "true"
	if e != nil {
		return ctx.Bad("No code.")
	}
	rt := ctx.FormValue("runtime")
	id := ctx.Param("problem")
	var problem contest.Problem
	ext := getExt(code.Filename)
	f, e := code.Open()
	if e != nil {
		return ctx.Bad("Could not open uploaded code!")
	}
	if ext == "" {
		return ctx.Bad("This file is not supported!")
	}
	if db.Database.NewSelect().Model(&problem).Where("id = ?", id).Scan(ctx.Request().Context()) != nil {
		return ctx.NotFound("Problem not found.")
	}
	if len(problem.Constraints.AllowedRuntimes) > 0 && !slices.Contains(problem.Constraints.AllowedRuntimes, rt) {
		return ctx.Bad("This language is not allowed by current problem.")
	}
	if !judge.Worker.RuntimeSupported(rt) {
		return ctx.InternalServerError("No judge server is available to handle this submission.")
	}
	dbSub, rollback, commit := createSubmission(ctx.Request().Context(), ctx.GetUUID(), problem.ID, rt, ext)
	if dbSub == nil {
		return ctx.Bad("Failed to create submission!")
	}
	p := storage.Submission.Create(dbSub.ID, ext)
	sub := prepare(dbSub.ID, p, rt, dbSub.AuthorID, &problem)
	if res, element, e = judge.Worker.Enqueue(sub, shouldStream, p, f); e != nil {
		judge.Worker.DestroySubscribers(sub.ID)
		rollback()
		return ctx.InternalServerError("Failed to process your submission.")
	}
	if commit() != nil {
		return ctx.InternalServerError("Could not save submission to database.")
	}
	if shouldStream && res != nil {
		stream := ctx.StreamResponse()
		go func() {
			<-ctx.Request().Context().Done()
			judge.Worker.Unsubscribe(sub.ID, element)
		}()
		for r := range res {
			stream.Write(r)
		}
		return ctx.Success()
	} else {
		return ctx.Respond(dbSub.ID)
	}
}
