// TODO: complete submit route

package problems

import (
	"context"
	"github.com/ArcticOJ/blizzard/v0/cache/stores"
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/models/contest"
	"github.com/ArcticOJ/blizzard/v0/judge"
	"github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/ArcticOJ/blizzard/v0/storage"
	"github.com/ArcticOJ/blizzard/v0/utils"
	"github.com/google/uuid"
	"path"
	"strings"
)

func prepare(id uint32, _path, language string, problem *contest.Problem) *judge.Submission {
	return &judge.Submission{
		ID:            id,
		Language:      language,
		SourcePath:    path.Base(_path),
		ProblemID:     problem.ID,
		TestCount:     problem.TestCount,
		PointsPerTest: problem.PointsPerTest,
		Constraints:   *problem.Constraints,
	}
}

func getExt(fileName string) string {
	return strings.ToLower(strings.TrimLeft(path.Ext(fileName), "."))
}

func createSubmission(ctx context.Context, userId uuid.UUID, problem, language string, ext string) (*contest.Submission, func() error, func() error) {
	sub := &contest.Submission{
		AuthorID:  userId,
		ProblemID: problem,
		Language:  language,
		Extension: ext,
	}
	tx, e := db.Database.Begin()
	if e != nil {
		return nil, nil, nil
	}
	if _, e = tx.NewInsert().Model(sub).Returning("id, submitted_at").Exec(ctx); e != nil {
		tx.Rollback()
		return nil, nil, nil
	}
	return sub, tx.Rollback, tx.Commit
}

func Submit(ctx *http.Context) http.Response {
	if ctx.RequireAuth() {
		return nil
	}
	code, e := ctx.FormFile("code")
	shouldStream := ctx.FormValue("stream") == "true"
	if e != nil {
		return ctx.Bad("No code.")
	}
	lang := ctx.FormValue("language")
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
	if len(problem.Constraints.AllowedLanguages) > 0 && !utils.ArrayIncludes(problem.Constraints.AllowedLanguages, lang) {
		return ctx.Bad("This language is not allowed by current problem.")
	}
	// might not be accurate
	// TODO: group judges by supported runtimes and do icmp pings on demand
	if !stores.Judge.IsRuntimeAllowed(ctx.Request().Context(), lang) {
		return ctx.InternalServerError("No judge server is available to handle this submission.")
	}
	dbSub, rollback, commit := createSubmission(ctx.Request().Context(), ctx.GetUUID(), problem.ID, lang, ext)
	if dbSub == nil {
		return ctx.Bad("Failed to create submission!")
	}
	p := storage.Submission.Create(dbSub.ID, ext)
	sub := prepare(dbSub.ID, p, lang, &problem)
	var res chan interface{}
	if shouldStream {
		res = judge.Worker.Subscribe(sub.ID)
	}
	if judge.Worker.Enqueue(sub, *dbSub.SubmittedAt) != nil || storage.Submission.Write(p, f) != nil {
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
			judge.Worker.Unsubscribe(sub.ID, res)
		}()
		for r := range res {
			stream.Write(r)
		}
		return ctx.Success()
	} else {
		return ctx.Respond(dbSub.ID)
	}
}
