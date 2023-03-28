package seed

import (
	"blizzard/blizzard/db/models/shared"
	"context"
	"github.com/uptrace/bun"
)

func registerModels(db *bun.DB) {
	db.RegisterModel((*shared.User)(nil))
}

func Populate(db *bun.DB) {
	registerModels(db)
	ctx := context.Background()
	if _, e := db.NewCreateTable().Model((*shared.User)(nil)).IfNotExists().Exec(ctx); e != nil {
		panic(e)
	}
}
