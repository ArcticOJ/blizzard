package cron

import (
	"context"
	"github.com/ArcticOJ/blizzard/v0/cron/jobs"
	"github.com/ArcticOJ/blizzard/v0/logger"
	"github.com/go-co-op/gocron/v2"
	"sync"
	"time"
)

var sched gocron.Scheduler

var once sync.Once

func init() {
	var e error
	sched, e = gocron.NewScheduler(gocron.WithLocation(time.UTC), gocron.WithStopTimeout(time.Second*3), gocron.WithLogger(Logger{}))
	logger.Panic(e, "failed to create cron scheduler")
}

func Start(ctx context.Context) {
	once.Do(func() {
		j, e := sched.NewJob(
			gocron.DurationJob(time.Minute*30),
			gocron.NewTask(jobs.PurgeSubmissions, ctx),
			gocron.WithName("purge-submissions"),
			gocron.WithSingletonMode(gocron.LimitModeReschedule),
			gocron.WithStartAt(gocron.WithStartImmediately()))
		logger.Panic(e, "failed to create cronjob '%s'", j.Name())
		sched.Start()
	})
}

func Stop() {
	sched.Shutdown()
}
