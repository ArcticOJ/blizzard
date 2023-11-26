package blizzard

import (
	"context"
	"github.com/ArcticOJ/blizzard/v0/cron"
	"github.com/ArcticOJ/blizzard/v0/judge"
)

func Init(ctx context.Context) {
	judge.Init(ctx)
	cron.Start(ctx)
}

func Destroy() {
	cron.Stop()
}
