package cache

import (
	"blizzard/blizzard/config"
	"blizzard/blizzard/logger"
	"context"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	cachestore "github.com/eko/gocache/lib/v4/store"
	bcstore "github.com/eko/gocache/store/bigcache/v4"
	"strings"
	"time"
)

type Store struct {
	ctx context.Context
	*marshaler.Marshaler
}

var store *Store

func Init() {
	ctx := context.Background()
	bc, _ := bigcache.NewBigCache(bigcache.DefaultConfig(5 * time.Minute))
	bcStore := bcstore.NewBigcache(bc)
	cacheManager := cache.New[interface{}](bcStore)
	store = &Store{Marshaler: marshaler.New(cacheManager), ctx: ctx}
}

func Get[T any](category, key string, loader func() (T, error)) *T {
	if config.Config.Debug {
		v, e := loader()
		if e != nil {
			return nil
		}
		return &v
	}
	if strings.TrimSpace(key) == "" {
		return nil
	}
	val, e := store.Get(store.ctx, key, new(T))
	if e != nil {
		v, e := loader()
		if e != nil {
			return nil
		}
		e = store.Set(store.ctx, key, v, cachestore.WithExpiration(time.Hour*24))
		return &v
	}
	if v, ok := val.(*T); ok {
		logger.Logger.Debug().Str("key", category+"_"+key).Interface("data", v).Msg("cache hit")
		return v
	} else {
		fmt.Println("fail")
	}
	return nil
}

func Invalidate(category, key string) bool {
	return store.Delete(store.ctx, category+"_"+key) == nil
}
