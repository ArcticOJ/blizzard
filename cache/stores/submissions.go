package stores

import (
	"github.com/ArcticOJ/blizzard/v0/cache"
	"github.com/ArcticOJ/blizzard/v0/rejson"
)

var Submissions *submissionStore

type submissionStore struct {
	j *rejson.ReJSON
}

const (
	// 30 minutes
	defaultExtraTtl = 30 * 60
)

func init() {
	Submissions = &submissionStore{j: &rejson.ReJSON{Client: cache.CreateClient(cache.Submission, "submissions")}}
}
