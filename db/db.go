package db

import (
	"blizzard/config"
	"blizzard/db/seed"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"net"
)

var Database *bun.DB

func createDb() *bun.DB {
	conf := config.Config.Database
	psql := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithUser(conf.Username),
		pgdriver.WithPassword(conf.Password),
		pgdriver.WithDatabase(conf.Name),
		pgdriver.WithAddr(net.JoinHostPort(conf.Host, fmt.Sprint(conf.Port))),
		pgdriver.WithInsecure(!conf.Secure)))
	db := bun.NewDB(psql, pgdialect.New())
	if config.Config.Debug {
		db.AddQueryHook(&DebugHook{})
	}
	return db
}

func init() {
	Database = createDb()
	seed.RegisterModels(Database)
}
