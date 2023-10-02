package debug

import (
	"blizzard/config"
	"blizzard/logger"
)

func Dump(obj ...interface{}) {
	if config.Config.Debug {
		logger.Blizzard.Debug().Interface("dump", obj).Msg("debug")
	}
}
