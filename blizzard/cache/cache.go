package cache

import (
	"blizzard/blizzard/config"
	"blizzard/blizzard/logger"
	"context"
	"fmt"
	"github.com/redis/rueidis"
	"time"
)

type Store struct {
	ctx context.Context
}

func init() {
}

func create(db int, name string) (c rueidis.Client) {
	c, e := rueidis.NewClient(rueidis.ClientOption{
		InitAddress:  []string{fmt.Sprintf("%s:%d", config.Config.Redis.Host, config.Config.Redis.Port)},
		SelectDB:     db,
		DisableCache: true,
	})
	logger.Panic(e, "failed to initialize redis client for %s cache", name)
	return
}

func Get[T any](key string, loader func() (T, error), duration time.Duration) *T {
	//if config.Config.Debug {
	//	v, e := loader()
	//	if e != nil {
	//		return nil
	//	}
	//	return &v
	//}
	//val, e := store.Get(store.ctx, key, new(T))
	//fmt.Println(e)
	//if e != nil {
	//	v, e := loader()
	//	if e != nil {
	//		return nil
	//	}
	//	e = store.Set(store.ctx, key, v, cachestore.WithExpiration(duration))
	//	return &v
	//}
	//if v, ok := val.(*T); ok {
	//	logger.Logger.Debug().Str("key", key).Interface("data", v).Msg("cache hit")
	//	return v
	//} else {
	//	Invalidate(key)
	//}
	//return nil
	v, e := loader()
	if e != nil {
		return nil
	}
	return &v
}
