// TODO: finalize submit route

package problems

import (
	"container/list"
	"context"
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/schema/contest"
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

func createSubmission(ctx context.Context, userId uuid.UUID, problem, runtime string, fileName string) (*contest.Submission, error) {
	sub := &contest.Submission{
		AuthorID:  userId,
		ProblemID: problem,
		Runtime:   runtime,
		FileName:  path.Base(fileName),
	}
	if _, e := db.Database.NewInsert().Model(sub).Returning("id, submitted_at, author_id").Exec(ctx); e != nil {
		return nil, e
	}
	return sub, nil
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
		return ctx.Bad("No source code.")
	}
	rt := ctx.FormValue("runtime")
	id := ctx.Param("problem")
	var problem contest.Problem
	// get file extension
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
		return ctx.Bad("This runtime is not allowed by current problem.")
	}
	if !judge.Observer.RuntimeSupported(rt) {
		return ctx.InternalServerError("No judge server is available to handle this submission.")
	}
	// create a random file name for this submission
	fName := storage.Submission.Create(ext)
	// try writing source code to previously generated file
	e, rollback := storage.Submission.Write(fName, f)
	if e != nil {
		return ctx.InternalServerError("Failed to write uploaded source code to disk.")
	}
	// write to database
	dbSub, _ := createSubmission(ctx.Request().Context(), ctx.GetUUID(), problem.ID, rt, fName)
	if dbSub == nil {
		rollback()
		return ctx.InternalServerError("Failed to commit current submission to database.")
	}
	// convert submission model to a polar-compatible one
	sub := prepare(dbSub.ID, fName, rt, dbSub.AuthorID, &problem)
	if res, element, e = judge.Observer.Enqueue(sub, shouldStream); e != nil {
		rollback()
		judge.Observer.DestroySubscribers(sub.ID)
		return ctx.InternalServerError("Failed to process your submission.")
	}
	if shouldStream && res != nil {
		stream := ctx.StreamResponse()
		done := ctx.Request().Context().Done()
		for {
			select {
			case <-done:
				judge.Observer.Unsubscribe(sub.ID, element)
				return nil
			case r, more := <-res:
				if !more || stream.Write(r) != nil {
					return nil
				}
			}
		}
	} else {
		return ctx.Respond(dbSub.ID)
	}
}
