package storage

import (
	"github.com/ArcticOJ/blizzard/v0/config"
	"github.com/ArcticOJ/blizzard/v0/logger"
	"os"
)

var (
	Submission submissionStorage
	READMEs    readmesStorage
)

func init() {
	for name, s := range config.Config.Blizzard.Storage {
		if e := os.Mkdir(s, 0755); e != nil && !os.IsExist(e) {
			logger.Blizzard.Fatal().Err(e).Msgf("error whilst creating storage for '%s'", name)
		}
	}
}
