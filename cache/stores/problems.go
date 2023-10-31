package stores

import (
	"context"
	"fmt"
	"github.com/ArcticOJ/blizzard/v0/cache"
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/models/contest"
	"github.com/ArcticOJ/blizzard/v0/rejson"
	"time"
)

var Problems *problemStore

type problemStore struct {
	j *rejson.ReJSON
}

const (
	defaultProblemKey = "blizzard::problem[%s]"
)

func init() {
	Problems = &problemStore{j: &rejson.ReJSON{Client: cache.CreateClient(cache.Problem, "problems")}}
}

func (s *problemStore) fallback(ctx context.Context, id string) *contest.Problem {
	prob := new(contest.Problem)
	if db.Database.NewSelect().Model(prob).Where("id = ?", id).Scan(ctx) != nil {
		return nil
	}
	// TODO: use global context to avoid interrupted operations
	s.j.JTxPipelined(ctx, func(r *rejson.ReJSON) error {
		k := fmt.Sprintf(defaultProblemKey, id)
		if e := r.JSONSet(ctx, k, "$", prob); e != nil {
			return e
		}
		return r.Expire(ctx, k, time.Hour*48).Err()
	})
	return prob
}

func (s *problemStore) Get(ctx context.Context, id string) *contest.Problem {
	p := s.j.JSONGet(ctx, fmt.Sprintf(defaultProblemKey, id))
	if r := rejson.Unmarshal[contest.Problem](p); r != nil {
		return r
	}
	return s.fallback(ctx, id)
}
