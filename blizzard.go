package blizzard

import (
	"context"
	"github.com/ArcticOJ/blizzard/v0/cron"
	"github.com/ArcticOJ/blizzard/v0/judge"
)

func Init(ctx context.Context) {
	judge.Init(ctx)
}

func Destroy(_ context.Context) {
	cron.Stop()
}
