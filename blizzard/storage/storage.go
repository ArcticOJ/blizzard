package storage

import (
	"blizzard/blizzard/config"
	"blizzard/blizzard/logger"
	"os"
)

var Submission SubmissionStorage

func init() {
	_s := config.Config.Storage
	for _, s := range []string{_s.Posts, _s.READMEs, _s.Problems, _s.Submissions} {
		if e := os.Mkdir(s, 0755); e != nil && !os.IsExist(e) {
			logger.Logger.Fatal().Err(e).Msgf("error whilst creating %d", s)
		}
	}
}
