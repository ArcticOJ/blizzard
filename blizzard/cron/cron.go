package cron

import (
	"blizzard/blizzard/cron/jobs"
	"github.com/go-co-op/gocron"
	"time"
)

func init() {
	s := gocron.NewScheduler(time.UTC)
	s.Every("10s").Do(jobs.UpdateJudgeStatus)
	s.StartAsync()
}
