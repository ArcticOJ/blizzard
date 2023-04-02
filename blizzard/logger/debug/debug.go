package debug

import (
	"blizzard/blizzard/config"
	"blizzard/blizzard/logger"
)

func Dump(obj ...interface{}) {
	if config.Config.Debug {
		logger.Logger.Debug().Interface("dump", obj).Msg("debug")
	}
}
