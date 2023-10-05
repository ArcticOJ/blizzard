package cron

import (
	"blizzard/cron/jobs"
	"context"
	"github.com/go-co-op/gocron"
	"sync"
	"time"
)

var sched *gocron.Scheduler

var once sync.Once

func init() {
	sched = gocron.NewScheduler(time.UTC)
	sched.SingletonModeAll()
}

func Start(ctx context.Context) {
	once.Do(func() {
		sched.Every("3s").Do(jobs.UpdateJudgeStatus, ctx)
		sched.Every("30m").Do(jobs.PurgeSubmissions, ctx)
		sched.StartAsync()
	})
}

func Stop() {
	sched.Stop()
}
