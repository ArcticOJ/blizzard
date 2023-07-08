package seed

import (
	"blizzard/blizzard/db/models/contest"
	"blizzard/blizzard/db/models/user"
	"context"
	"github.com/uptrace/bun"
)

var intermediaryModels = []any{
	(*contest.ContestToOrganizer)(nil),
	(*contest.ContestToProblem)(nil),
	(*user.UserToRole)(nil),
}

var models = []any{
	// users
	(*user.User)(nil),
	(*user.OAuthConnection)(nil),
	(*user.Role)(nil),
	(*contest.Contest)(nil),
	(*contest.Problem)(nil),
	(*contest.Submission)(nil),
	// intermediary models
}

func RegisterModels(db *bun.DB) {
	db.RegisterModel(intermediaryModels...)
	db.RegisterModel(models...)
}

func InitTables(db *bun.DB, ctx context.Context) error {
	if e := db.ResetModel(ctx, models...); e != nil {
		return e
	}
	return db.ResetModel(ctx, intermediaryModels...)
}
