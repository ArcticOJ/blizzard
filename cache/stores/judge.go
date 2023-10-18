package stores

import (
	"context"
	"github.com/ArcticOJ/blizzard/v0/cache"
	"github.com/redis/go-redis/v9"
	"time"
)

var Judge *judgeStatusStore

type (
	JudgeStatus struct {
		Version     string         `json:"version"`
		Memory      uint32         `json:"memory"`
		OS          string         `json:"os"`
		Parallelism uint8          `json:"parallelism"`
		BootedSince uint64         `json:"bootedSince"`
		Runtimes    []JudgeRuntime `json:"runtimes"`
	}

	JudgeRuntime struct {
		ID        string `json:"id"`
		Compiler  string `json:"compiler"`
		Arguments string `json:"arguments"`
		Version   string `json:"version"`
	}

	judgeStatusStore struct {
		c redis.UniversalClient
	}
)

const (
	defaultJudgeStatusKey     = "blizzard::judge_status"
	defaultAllowedRuntimesKey = "blizzard::allowed_runtimes"
	defaultJudgeListKey       = "blizzard::judge_list"
)

func init() {
	Judge = &judgeStatusStore{c: cache.CreateClient(cache.Judge, "judge")}
}

func (s *judgeStatusStore) UpdateJudgeStatus(ctx context.Context, judgeList []interface{}, status string, allowedRuntimes []interface{}) {
	s.c.TxPipelined(ctx, func(p redis.Pipeliner) error {
		s.c.Set(ctx, defaultJudgeStatusKey, status, time.Hour*24)
		p.Del(ctx, defaultAllowedRuntimesKey).Err()
		p.SAdd(ctx, defaultAllowedRuntimesKey, allowedRuntimes...)
		p.SAdd(ctx, defaultJudgeListKey, judgeList...)
		return nil
	})
}

func (s *judgeStatusStore) IsRuntimeAllowed(ctx context.Context, runtime string) bool {
	v, e := s.c.SIsMember(ctx, defaultAllowedRuntimesKey, runtime).Result()
	return v && e == nil
}

func (s *judgeStatusStore) GetJudgeList(ctx context.Context) []string {
	v, e := s.c.SMembers(ctx, defaultJudgeListKey).Result()
	if e != nil {
		return nil
	}
	return v
}

func (s *judgeStatusStore) GetJudgeStatus(ctx context.Context) []byte {
	status, e := s.c.Get(ctx, defaultJudgeStatusKey).Result()
	if e != nil {
		return []byte("null")
	}
	return []byte(status)
}
