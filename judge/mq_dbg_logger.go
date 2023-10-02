package judge

import (
	"blizzard/config"
	"blizzard/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

func init() {
	if config.Config.Debug {
		amqp.SetLogger(DbgLogger{})
	}
}

type DbgLogger struct{}

func (DbgLogger) Printf(format string, v ...interface{}) {
	logger.Blizzard.Debug().Msgf(format, v)
}
