package cache

import (
	"context"
	"fmt"
	"github.com/ArcticOJ/blizzard/v0/config"
	"github.com/ArcticOJ/blizzard/v0/logger"
	"github.com/redis/go-redis/v9"
	"net"
	"time"
)

type (
	DB uint8
)

const (
	Result DB = iota + 1
	User
	Bucket
	Submission
	Judge
	Problem
)

func CreateClient(db DB, name string) (c redis.UniversalClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	c = redis.NewClient(&redis.Options{
		Addr: net.JoinHostPort(config.Config.Dragonfly.Host, fmt.Sprint(config.Config.Dragonfly.Port)),
		DB:   int(db),
	})
	if config.Config.Debug {
		c.AddHook(DebugHook{Name: name})
	}
	logger.Panic(c.Ping(ctx).Err(), "failed to initialize redis client for %s cache", name)
	return
}
