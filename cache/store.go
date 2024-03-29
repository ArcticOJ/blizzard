package cache

import (
	"fmt"
	"github.com/ArcticOJ/blizzard/v0/config"
	"github.com/redis/go-redis/v9"
	"net"
)

type (
	DB uint8
)

const (
	User DB = iota + 1
	Bucket
	Submission
	Problem
)

func CreateClient(db DB, name string) (c redis.UniversalClient) {
	c = redis.NewClient(&redis.Options{
		Addr: net.JoinHostPort(config.Config.Blizzard.Dragonfly.Host, fmt.Sprint(config.Config.Blizzard.Dragonfly.Port)),
		DB:   int(db),
	})
	if config.Config.Debug {
		c.AddHook(DebugHook{Name: name})
	}
	return
}
