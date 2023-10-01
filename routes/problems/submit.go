// TODO: complete submit route

package problems

import (
	"blizzard/core"
	"blizzard/core/errs"
	"blizzard/db"
	"blizzard/db/models/contest"
	"blizzard/judge"
	"blizzard/logger"
	"blizzard/models"
	"blizzard/models/extra"
	"blizzard/storage"
	"context"
	"errors"
	"github.com/google/uuid"
	"path"
	"strings"
)

func prepare(id uint32, _path, language string, problem *contest.Problem) *judge.Submission {
	return &judge.Submission{
		ID:          id,
		Language:    language,
		SourcePath:  path.Base(_path),
		ProblemID:   problem.ID,
		TestCount:   problem.TestCount,
		Constraints: problem.Constraints,
	}
}

func getExt(language, fileName string) string {
	ext := strings.ToLower(strings.TrimLeft(path.Ext(fileName), "."))
	if ext == "" {
		if l, ok := core.LanguageMatrix[language]; ok {
			return l.Extension
		}
		return ""
	}
	for _, l := range core.LanguageMatrix {
		if l.Extension == ext {
			return ext
		}
	}
	return ""
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

// TODO: check availability of judges before judging

func Submit(ctx *extra.Context) models.Response {
	if ctx.RequireAuth() {
		return nil
	}
	code, e := ctx.FormFile("code")
	if e != nil {
		return ctx.Bad("No code.")
	}
	lang := ctx.FormValue("language")
	logger.Logger.Debug().Str("lang", lang).Send()
	id := ctx.Param("problem")
	var problem contest.Problem
	if _, ok := core.LanguageMatrix[lang]; !ok {
		return ctx.Bad("Unsupported language!")
	}
	ext := getExt(lang, code.Filename)
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
	if !judge.ResponseObserver.CheckAvailability(lang, ctx.Request().Context()) {
		return ctx.InternalServerError("No judge server is available to handle this submission.")
	}
	dbSub, rollback, commit := createSubmission(ctx.Request().Context(), *ctx.GetUUID(), problem.ID, lang, ext)
	if dbSub == nil {
		return ctx.Bad("Failed to create submission!")
	}
	p := storage.Submission.Create(dbSub.ID, ext)
	sub := prepare(dbSub.ID, p, lang, &problem)
	if storage.Submission.Write(p, f) != nil {
		return ctx.Bad("Could not write code to file!")
	}
	res := make(chan interface{}, 1)
	s := judge.ResponseObserver.Observe(sub.ID, res)
	if judge.ResponseObserver.Enqueue(sub, *dbSub.SubmittedAt) != nil {
		judge.ResponseObserver.DestroyObserver(sub.ID)
		rollback()
		if errors.Is(e, errs.JudgeNotAvailable) {
			return ctx.InternalServerError("No judge is available for this language.")
		}
		return ctx.InternalServerError("Failed to enqueue submission.")
	}
	commit()
	stream := ctx.StreamResponse()
	go func() {
		<-ctx.Request().Context().Done()
		judge.ResponseObserver.StopObserve(s)
	}()
	for r := range res {
		stream.Write(r)
	}
	return ctx.Success()
}
