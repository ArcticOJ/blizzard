package migrations

import (
	"context"
	"github.com/ArcticOJ/blizzard/v0/db/seed"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		return seed.InitTables(db, ctx)
	}, nil)
}
