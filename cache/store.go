package cache

import (
	"blizzard/config"
	"blizzard/logger"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net"
	"time"
)

type (
	DB uint8
)

const (
	Result DB = iota + 1
	Submission
	User
	Bucket
	Judge
)

func CreateClient(db DB, name string) (c *redis.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	c = redis.NewClient(&redis.Options{
		Addr: net.JoinHostPort(config.Config.Redis.Host, fmt.Sprint(config.Config.Redis.Port)),
		DB:   int(db),
	})
	if config.Config.Debug {
		c.AddHook(DebugHook{Name: name})
	}
	logger.Panic(c.Ping(ctx).Err(), "failed to initialize redis client for %s cache", name)
	return
}
