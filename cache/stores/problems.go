package stores

import (
	"github.com/ArcticOJ/blizzard/v0/cache"
	"github.com/ArcticOJ/blizzard/v0/rejson"
)

var Problems *ProblemStore

type ProblemStore struct {
	j *rejson.ReJSON
}

const (
	defaultProblemKey = "blizzard::problem[%s]"
)

func init() {
	Problems = &ProblemStore{j: &rejson.ReJSON{Client: cache.CreateClient(cache.Problem, "problems")}}
}

func (*ProblemStore) Get() {

}
