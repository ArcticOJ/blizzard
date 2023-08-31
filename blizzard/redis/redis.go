package redis

import (
	"blizzard/blizzard/config"
	"blizzard/blizzard/logger"
	"fmt"
	"github.com/redis/rueidis"
)

var CacheClient rueidis.Client
var LockerRedis rueidis.Client

func init() {
	var e error
	CacheClient, e = rueidis.NewClient(rueidis.ClientOption{
		InitAddress:  []string{fmt.Sprintf("%s:%d", config.Config.Redis.Host, config.Config.Redis.Port)},
		SelectDB:     config.Config.Redis.DB,
		DisableCache: true,
	})
	logger.Panic(e, "failed to initialize redis client for cache storage")
	LockerRedis, e = rueidis.NewClient(rueidis.ClientOption{
		InitAddress:  []string{fmt.Sprintf("%s:%d", config.Config.Redis.Host, config.Config.Redis.Port)},
		SelectDB:     config.Config.Redis.DB + 1,
		DisableCache: true,
	})
	logger.Panic(e, "failed to initialize redis client for distributed locker")
}
