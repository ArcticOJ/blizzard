package blizzard

import (
	"blizzard/cron"
	"blizzard/judge"
	"context"
)

func Init(ctx context.Context) {
	judge.ResponseObserver = judge.NewObserver(ctx)
	cron.Start(ctx)
	go judge.ResponseObserver.Work()
}

func Destroy() {
	cron.Stop()
}
