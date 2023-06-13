package seed

import (
	"blizzard/blizzard/db/models/users"
	"context"
	"github.com/uptrace/bun"
)

var models = []any{
	(*users.User)(nil),
	(*users.OAuthConnection)(nil),
}

func registerModels(db *bun.DB) {
	for _, model := range models {
		db.RegisterModel(model)
	}
}

func Populate(db *bun.DB) {
	registerModels(db)
	ctx := context.Background()
	for _, model := range models {
		if _, e := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx); e != nil {
			panic(e)
		}
	}
}
