package blizzard

import (
	"context"
	"github.com/ArcticOJ/blizzard/v0/cron"
	"github.com/ArcticOJ/blizzard/v0/judge"
)

func Init(ctx context.Context) {
	judge.ResponseWorker = judge.NewWorker(ctx)
	cron.Start(ctx)
	go judge.ResponseWorker.Work()
}

func Destroy() {
	cron.Stop()
}
