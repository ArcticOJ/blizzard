package stores

import (
	"context"
	"fmt"
	"github.com/ArcticOJ/blizzard/v0/cache"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var Pending *PendingStore

type PendingStore struct {
	c  redis.UniversalClient
	l  sync.RWMutex
	tl sync.RWMutex
}

const (
	defaultPendingSubmissionKey = "blizzard::pending_submission[%d]"
	// 30 minutes
	defaultExtraTtl = 30 * 30
)

func init() {
	Pending = &PendingStore{c: cache.CreateClient(cache.Result, "results")}
}

func (s *PendingStore) Set(ctx context.Context, id uint32, tag uint64, count, ttl uint16) error {
	s.l.Lock()
	defer s.l.Unlock()
	return s.c.SetEx(ctx, fmt.Sprintf(defaultPendingSubmissionKey, id), tag, time.Duration(defaultExtraTtl+count*(ttl+2))*time.Second).Err()
}

func (s *PendingStore) IsPending(ctx context.Context, id uint32) bool {
	s.l.RLock()
	defer s.l.RUnlock()
	ok, e := s.c.Exists(ctx, fmt.Sprintf(defaultPendingSubmissionKey, id)).Result()
	return ok == 1 && e == nil
}

func (s *PendingStore) Get(ctx context.Context, id uint32) (uint64, bool) {
	tag, e := s.c.Get(ctx, fmt.Sprintf(defaultPendingSubmissionKey, id)).Uint64()
	return tag, e == nil
}

func (s *PendingStore) Delete(ctx context.Context, id uint32) {
	s.l.Lock()
	defer s.l.Unlock()
	s.c.Del(ctx, fmt.Sprintf(defaultPendingSubmissionKey, id))
}
