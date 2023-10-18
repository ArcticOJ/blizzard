package stores

import (
	"context"
	"fmt"
	"github.com/ArcticOJ/blizzard/v0/cache"
	"github.com/ArcticOJ/blizzard/v0/db/models/contest"
	"github.com/ArcticOJ/blizzard/v0/rejson"
	"github.com/ArcticOJ/blizzard/v0/utils"
	"sync"
	"time"
)

var Submissions *submissionStore

type submissionStore struct {
	j *rejson.ReJSON
	l sync.RWMutex
}

const (
	defaultPendingSubmissionKey = "blizzard::pending_submission[%d]"
	// 30 minutes
	defaultExtraTtl = 30 * 60
)

func init() {
	Submissions = &submissionStore{j: &rejson.ReJSON{Client: cache.CreateClient(cache.Submission, "submissions")}}
}

func (s *submissionStore) SetPending(ctx context.Context, id uint32, tag uint64, count, ttl uint16) error {
	s.l.Lock()
	defer s.l.Unlock()
	return s.j.JTxPipelined(ctx, func(r *rejson.ReJSON) error {
		k := fmt.Sprintf(defaultPendingSubmissionKey, id)
		if e := r.JSONSet(ctx, k, "$", map[string]interface{}{
			"tag":   tag,
			"cases": utils.ArrayFill[interface{}](nil, int(count)),
		}); e != nil {
			return e
		}
		return r.Expire(ctx, k, time.Duration(defaultExtraTtl+count*(ttl+2))*time.Second).Err()
	})
}

func (s *submissionStore) UpdatePending(ctx context.Context, id uint32, result contest.CaseResult) error {
	return s.j.JSONSet(ctx, fmt.Sprintf(defaultPendingSubmissionKey, id), fmt.Sprintf("$..cases[%d]", result.ID-1), result)
}

func (s *submissionStore) IsPending(ctx context.Context, id uint32) bool {
	s.l.RLock()
	defer s.l.RUnlock()
	ok, e := s.j.Exists(ctx, fmt.Sprintf(defaultPendingSubmissionKey, id)).Result()
	return ok == 1 && e == nil
}

func (s *submissionStore) GetPendingResults(ctx context.Context, id uint32) []contest.CaseResult {
	r := s.j.JSONGet(ctx, fmt.Sprintf(defaultPendingSubmissionKey, id), "$.cases")
	if _r := rejson.Unmarshal[[][]contest.CaseResult](r); _r != nil && len(*_r) > 0 {
		return (*_r)[0]
	}
	return nil
}

func (s *submissionStore) GetPendingTag(ctx context.Context, id uint32) (uint64, bool) {
	r := s.j.JSONGet(ctx, fmt.Sprintf(defaultPendingSubmissionKey, id), "$")
	if t := rejson.Unmarshal[[]uint64](r); t != nil && len(*t) > 0 {
		return (*t)[0], true
	}
	return 0, false
}

func (s *submissionStore) DeletePending(ctx context.Context, id uint32) {
	s.l.Lock()
	defer s.l.Unlock()
	s.j.Del(ctx, fmt.Sprintf(defaultPendingSubmissionKey, id))
}
