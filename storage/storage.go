package storage

import (
	"blizzard/config"
	"blizzard/logger"
	"os"
)

var Submission SubmissionStorage

func init() {
	_s := config.Config.Storage
	for _, s := range []string{_s.Posts, _s.READMEs, _s.Problems, _s.Submissions} {
		if e := os.Mkdir(s, 0755); e != nil && !os.IsExist(e) {
			logger.Blizzard.Fatal().Err(e).Msgf("error whilst creating %d", s)
		}
	}
}
