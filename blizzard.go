package blizzard

import (
	"blizzard/cron"
	"blizzard/judge"
	"context"
)

func Init(ctx context.Context) {
	judge.ResponseWorker = judge.NewObserver(ctx)
	cron.Start(ctx)
	go judge.ResponseWorker.Work()
}

func Destroy() {
	cron.Stop()
}
