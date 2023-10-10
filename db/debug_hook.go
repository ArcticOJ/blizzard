package db

import (
	"context"
	"errors"
	"github.com/ArcticOJ/blizzard/v0/logger"
	"github.com/uptrace/bun"
	"reflect"
	"time"
)

type DebugHook struct {
}

var _ bun.QueryHook = (*DebugHook)(nil)

func (*DebugHook) Init(_ *bun.DB) {
	logger.Blizzard.Debug().Msg("initializing database")
}

func (*DebugHook) BeforeQuery(ctx context.Context, _ *bun.QueryEvent) context.Context {
	return ctx
}

func (*DebugHook) AfterQuery(_ context.Context, event *bun.QueryEvent) {
	now := time.Now()
	dur := now.Sub(event.StartTime)
	log := logger.Blizzard.Debug().Time("timestamp", now).Str("operation", event.Operation()).Str("duration", dur.Round(time.Microsecond).String()).Str("query", event.Query)
	if event.Err != nil {
		typ := reflect.TypeOf(event.Err).String()
		log = log.Err(errors.New(typ + ": " + event.Err.Error()))
	}
	log.Msg("bun")
}
