package stores

import (
	"blizzard/cache"
	"blizzard/db"
	"blizzard/db/models/user"
	"blizzard/logger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/nitishm/go-rejson/v4"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
	"time"
)

var Users *UserStore

type UserStore struct {
	c *redis.Client
	j *rejson.Handler
}

const (
	defaultUserListKey        = "blizzard::user_list"
	defaultUserKey            = "blizzard::user[%s]"
	defaultHandleToIdResolver = "blizzard::user_id[%s]"
)

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	Users = &UserStore{cache.CreateClient(cache.User, "users"), rejson.NewReJSONHandler()}
	Users.j.SetGoRedisClient(Users.c)
	var ids []string
	logger.Panic(db.Database.NewSelect().Model((*user.User)(nil)).Column("id").Scan(ctx, &ids), "failed to to query for users")
	var m []redis.Z
	for _, id := range ids {
		m = append(m, redis.Z{
			Score:  0,
			Member: id,
		})
	}
	_, e := Users.c.TxPipelined(ctx, func(p redis.Pipeliner) error {
		if e := p.Del(ctx, defaultUserListKey).Err(); e != nil {
			return e
		}
		return p.ZAdd(ctx, defaultUserListKey, m...).Err()
	})
	logger.Panic(e, "failed to populate user cache")
}

func (c *UserStore) load(id uuid.UUID, handle string) (u *user.User, _id uuid.UUID) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	u = new(user.User)
	q := db.Database.NewSelect().Model(u)
	if id == uuid.Nil {
		q = q.Where("handle = ?", handle)
	} else {
		q = q.Where("id = ?", id)
	}
	if e := q.Relation("Connections", func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Where("show_in_profile = true")
	}).Relation("Roles", func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Order("priority ASC").Column("name", "icon", "color")
	}).Scan(ctx); e != nil {
		return nil, uuid.Nil
	}
	if len(u.Roles) > 0 {
		u.TopRole = &u.Roles[0]
	}
	if u.ID != uuid.Nil {
		_id = u.ID
		handle = u.Handle
		if e := c.c.Set(ctx, fmt.Sprintf(defaultHandleToIdResolver, handle), u.ID.String(), time.Hour*6).Err(); e != nil {
			logger.Logger.Error().Err(e).Str("id", id.String()).Str("handle", handle).Msgf("could not bind handle to id")
		}
	}
	return
}

func (c *UserStore) Exists(ctx context.Context, id uuid.UUID) bool {
	exists, e := c.c.Exists(ctx, fmt.Sprintf(defaultUserKey, id)).Result()
	return exists == 0 && e == nil
}

func (c *UserStore) fallback(ctx context.Context, id uuid.UUID, handle string) *user.User {
	u, _id := c.load(id, handle)
	if u != nil && _id != uuid.Nil {
		c.j.JSONSet(ctx, fmt.Sprintf(defaultUserKey, _id), "$", u)
	}
	return u
}

func (c *UserStore) Get(ctx context.Context, id uuid.UUID, handle string) *user.User {
	if handle == "" && id == uuid.Nil {
		return nil
	}
	if id == uuid.Nil {
		if _id, e := c.c.Get(ctx, fmt.Sprintf(defaultHandleToIdResolver, handle)).Result(); e == nil && _id != "" && _id != uuid.Nil.String() {
			id, e = uuid.Parse(_id)
			if e == nil {
				return c.fallback(ctx, id, "")
			}
		}
		return c.fallback(ctx, uuid.Nil, handle)
	}
	if c.c.ZScore(ctx, defaultUserListKey, id.String()).Err() == redis.Nil {
		return nil
	}
	r, e := c.j.JSONGet(ctx, fmt.Sprintf(defaultUserKey, id), "$")
	if e == redis.Nil {
		return c.fallback(ctx, id, "")
	}
	var _u []user.User
	if json.Unmarshal(r.([]byte), &_u) != nil {
		return c.fallback(ctx, id, "")
	} else if len(_u) > 0 {
		return &_u[0]
	}
	return nil
}
