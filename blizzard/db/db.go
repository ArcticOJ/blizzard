package db

import (
	"blizzard/blizzard/config"
	"blizzard/blizzard/db/seed"
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var Database *bun.DB

func createDb() *bun.DB {
	conf := config.Config.Database
	psql := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithUser(conf.Username),
		pgdriver.WithPassword(conf.Password),
		pgdriver.WithDatabase(conf.Name),
		pgdriver.WithAddr(conf.Address),
		pgdriver.WithInsecure(!conf.Secure)))
	db := bun.NewDB(psql, pgdialect.New())
	if config.Config.Debug {
		db.AddQueryHook(&QueryHook{})
	}
	return db
}

func Init() {
	Database = createDb()
	seed.RegisterModels(Database)
}
