package cron

import "github.com/ArcticOJ/blizzard/v0/logger"

type Logger struct{}

func (c Logger) Debug(msg string, args ...interface{}) {
	logger.Blizzard.Debug().Any("args", args).Msgf("cron: %s", msg)
}

func (c Logger) Error(msg string, args ...interface{}) {
	logger.Blizzard.Error().Any("args", args).Msgf("cron: %s", msg)
}

func (c Logger) Info(msg string, args ...interface{}) {
	logger.Blizzard.Info().Any("args", args).Msgf("cron: %s", msg)
}

func (c Logger) Warn(msg string, args ...interface{}) {
	logger.Blizzard.Warn().Any("args", args).Msgf("cron: %s", msg)
}
