package debug

import (
	"github.com/ArcticOJ/blizzard/v0/config"
	"github.com/ArcticOJ/blizzard/v0/logger"
)

func Dump(obj ...interface{}) {
	if config.Config.Debug {
		logger.Blizzard.Debug().Interface("dump", obj).Msg("debug")
	}
}
