package db

import (
	"blizzard/blizzard/logger"
	"context"
	"errors"
	"github.com/uptrace/bun"
	"reflect"
	"time"
)

type QueryHook struct {
}

var _ bun.QueryHook = (*QueryHook)(nil)

func (*QueryHook) Init(_ *bun.DB) {
	logger.Logger.Debug().Msg("initializing database")
}

func (*QueryHook) BeforeQuery(ctx context.Context, _ *bun.QueryEvent) context.Context {
	return ctx
}

func (*QueryHook) AfterQuery(_ context.Context, event *bun.QueryEvent) {
	now := time.Now()
	dur := now.Sub(event.StartTime)
	log := logger.Logger.Debug().Time("timestamp", now).Str("operation", event.Operation()).Str("duration", dur.Round(time.Microsecond).String()).Str("query", event.Query)
	if event.Err != nil {
		typ := reflect.TypeOf(event.Err).String()
		log = log.Err(errors.New(typ + ": " + event.Err.Error()))
	}
	log.Msg("bun")
}
