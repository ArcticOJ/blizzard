package stores

import (
	"blizzard/cache"
	"github.com/redis/go-redis/v9"
)

var Session *SessionStore

const defaultSessionKey = "blizzard::session[%s]"

type SessionStore struct {
	c *redis.Client
}

func init() {
	Session = &SessionStore{cache.CreateClient(cache.Session, "sessions")}
}
