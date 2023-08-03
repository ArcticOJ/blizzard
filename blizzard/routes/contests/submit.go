// TODO: complete submit route

package contests

import (
	"blizzard/blizzard/db"
	"blizzard/blizzard/db/models/contest"
	"blizzard/blizzard/judge"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"blizzard/blizzard/pb"
	"blizzard/blizzard/utils"
	"context"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
)

func prepare(id uint32, language string, problem *contest.Problem, file *multipart.FileHeader) *pb.Submission {
	f, e := file.Open()
	if e != nil {
		return nil
	}
	buf, e := io.ReadAll(f)
	if e != nil {
		return nil
	}
	c := problem.Constraints
	return &pb.Submission{
		Id:       id,
		Language: language,
		Buffer:   buf,
		Problem:  problem.ID,
		Metadata: &pb.Metadata{
			ShortCircuit: c.ShortCircuit,
			CaseCount:    uint32(problem.TestCount),
			Duration:     c.TimeLimit,
			Memory:       c.MemoryLimit,
		},
		Checker: utils.ReadFile("/data/Dev/mock/test.py"),
	}
}

var id uint32 = 0

func createSubmission(ctx context.Context, userId uuid.UUID, problem string, language string) *contest.Submission {
	id++
	sub := &contest.Submission{
		ID:        id,
		AuthorID:  userId,
		ProblemID: problem,
		Language:  language,
	}
	/*if _, e := db.Database.NewInsert().Model(sub).Returning("id").Exec(ctx); e != nil {
		return nil
	}*/
	return sub
}

func Submit(ctx *extra.Context) models.Response {
	if ctx.RequireAuth() {
		return nil
	}
	code, e := ctx.FormFile("code")
	if e != nil {
		return ctx.Bad("No code.")
	}
	// TODO: check if language is valid
	lang := ctx.FormValue("language")
	id := ctx.Param("id")
	var problem contest.Problem
	if db.Database.NewSelect().Model(&problem).Where("id = ?", id).Scan(ctx.Request().Context()) != nil {
		return ctx.NotFound("Problem not found.")
	}
	dbSub := createSubmission(ctx.Request().Context(), *ctx.GetUUID(), problem.ID, lang)
	if dbSub == nil {
		return ctx.InternalServerError("Failed to create submission!")
	}
	sub := prepare(dbSub.ID, lang, &problem, code)
	if _judge := judge.Enqueue(sub, dbSub); _judge != nil {
		// number of cases + final response
		res := make(chan interface{}, sub.Metadata.CaseCount+1)
		_judge()
		judge.ResultWatcher.Track(sub.Id, res)
		stream := ctx.StreamResponse()
		go func() {
			<-ctx.Request().Context().Done()
			judge.ResultWatcher.Untrack(sub.Id, res)
		}()
		for r := range res {
			stream.Write(r)
		}
	} else {
		return ctx.InternalServerError("Could not find a suitable judge server.")
	}
	return nil
}
