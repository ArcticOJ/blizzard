package stores

import (
	"context"
	"fmt"
	"github.com/ArcticOJ/blizzard/v0/cache"
	"github.com/ArcticOJ/blizzard/v0/config"
	"github.com/ArcticOJ/blizzard/v0/logger"
	"github.com/redis/go-redis/v9"
	"math"
)

var RateLimiter *rateLimitStore

const defaultBucketKey = "blizzard::bucket[%s]"

type rateLimitStore struct {
	c redis.UniversalClient
}

func init() {
	RateLimiter = &rateLimitStore{cache.CreateClient(cache.Bucket, "buckets")}
}

func (s *rateLimitStore) Limit(ctx context.Context, ip string) (allowed bool, totalLimit, remaining, retryAfter, nextReset int64) {
	v, e := s.c.Do(ctx, "CL.THROTTLE", fmt.Sprintf(defaultBucketKey, ip), uint16(math.Max(math.Ceil(float64(config.Config.Blizzard.RateLimit)/2), 1)), config.Config.Blizzard.RateLimit, 30, 1).Int64Slice()
	if e != nil || len(v) != 5 {
		logger.Blizzard.Err(e).Msgf("failed to process rate limit for '%s'", ip)
		return false, 0, 0, 0, 0
	}
	return v[0] == 0, v[1], v[2], v[3], v[4]
}
