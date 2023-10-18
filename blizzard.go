package blizzard

import (
	"context"
	"github.com/ArcticOJ/blizzard/v0/cron"
	"github.com/ArcticOJ/blizzard/v0/judge"
)

func Init(ctx context.Context) {
	judge.Worker = judge.NewWorker(ctx)
	cron.Start(ctx)
	go judge.Worker.Work()
}

func Destroy() {
	cron.Stop()
}