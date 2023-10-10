package seed

import (
	"context"
	"github.com/ArcticOJ/blizzard/v0/db/models/contest"
	"github.com/ArcticOJ/blizzard/v0/db/models/post"
	"github.com/ArcticOJ/blizzard/v0/db/models/user"
	"github.com/uptrace/bun"
)

var intermediaryModels = []any{
	(*contest.ContestToOrganizer)(nil),
	(*contest.ContestToProblem)(nil),
	(*user.UserToRole)(nil),
}

var models = []any{
	(*user.User)(nil),
	(*user.OAuthConnection)(nil),
	(*user.Role)(nil),

	(*contest.Contest)(nil),
	(*contest.Problem)(nil),
	(*contest.Submission)(nil),

	(*post.Post)(nil),
	(*post.Comment)(nil),
}

func RegisterModels(db *bun.DB) {
	db.RegisterModel(intermediaryModels...)
	db.RegisterModel(models...)
}

func DropAll(db *bun.DB, ctx context.Context) error {
	m := append(models, intermediaryModels...)
	for i := range m {
		if _, e := db.NewDropTable().Model(m[i]).Cascade().IfExists().Exec(ctx); e != nil {
			return e
		}
	}
	return nil
}

func CreateAll(db *bun.DB, ctx context.Context) error {
	m := append(models, intermediaryModels...)
	for i := range m {
		if _, e := db.NewCreateTable().Model(m[i]).Exec(ctx); e != nil {
			return e
		}
	}
	return nil
}

func InitTables(db *bun.DB, ctx context.Context) error {
	if e := db.ResetModel(ctx, models...); e != nil {
		return e
	}
	return db.ResetModel(ctx, intermediaryModels...)
}
