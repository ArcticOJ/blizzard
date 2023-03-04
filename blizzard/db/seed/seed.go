package seed

import (
	"backend/blizzard/db/models/shared"
	"context"
	"github.com/uptrace/bun"
)

func RegisterModels(db *bun.DB) {
	db.RegisterModel((*shared.User)(nil))
}

func Populate(db *bun.DB) {
	RegisterModels(db)
	ctx := context.Background()
	//db.NewRaw("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Scan(ctx)
	if _, e := db.NewCreateTable().Model((*shared.User)(nil)).IfNotExists().Exec(ctx); e != nil {
		panic(e)
	}
}
