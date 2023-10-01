package migrations

import (
	"blizzard/db/seed"
	"context"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		return seed.InitTables(db, ctx)
	}, nil)
}
