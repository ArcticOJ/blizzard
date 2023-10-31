package judge

import (
	"github.com/ArcticOJ/blizzard/v0/config"
	"github.com/ArcticOJ/blizzard/v0/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

func init() {
	if config.Config.Debug {
		amqp.SetLogger(dbgLogger{})
	}
}

type dbgLogger struct{}

func (dbgLogger) Printf(format string, v ...interface{}) {
	logger.Blizzard.Debug().Msgf(format, v)
}
