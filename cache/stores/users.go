package stores

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/ArcticOJ/blizzard/v0/cache"
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/models/user"
	"github.com/ArcticOJ/blizzard/v0/rejson"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/tmthrgd/go-hex"
	"github.com/uptrace/bun"
	"strings"
	"time"
)

var Users *UserStore

type UserStore struct {
	j *rejson.ReJSON
}

const (
	defaultUserKey            = "blizzard::user[%s]"
	defaultHandleToIdResolver = "blizzard::user_id[%s]"
)

func init() {
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	//defer cancel()
	Users = &UserStore{j: &rejson.ReJSON{Client: cache.CreateClient(cache.User, "users")}}
	//var ids []string
	//logger.Panic(db.Database.NewSelect().Model((*user.User)(nil)).Column("id").Scan(ctx, &ids), "failed to to query for users")
	//var m []redis.Z
	//for _, id := range ids {
	//	m = append(m, redis.Z{
	//		Score:  0,
	//		Member: id,
	//	})
	//}
	//_, e := Users.j.TxPipelined(ctx, func(p redis.Pipeliner) error {
	//	p.Del(ctx, defaultUserListKey)
	//	p.ZAdd(ctx, defaultUserListKey, m...)
	//	return nil
	//})
	//logger.Panic(e, "failed to populate user cache")
}

func (s *UserStore) load(id uuid.UUID, handle string) (u *user.User) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	u = new(user.User)
	q := db.Database.NewSelect().Model(u).ExcludeColumn("password", "api_key")
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
		return nil
	}
	if strings.TrimSpace(u.Email) != "" {
		h := md5.Sum([]byte(u.Email))
		u.Avatar = hex.EncodeToString(h[:])
	}
	return
}

//func (s *UserStore) Exists(ctx context.Context, id uuid.UUID) bool {
//	return s.j.ZScore(ctx, defaultUserListKey, id.String()).Err() == nil
//}

func (s *UserStore) fallback(ctx context.Context, id uuid.UUID, handle string) *user.User {
	u := s.load(id, handle)
	if u != nil && u.ID != uuid.Nil {
		handle = u.Handle
		s.j.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
			p := &rejson.ReJSON{Client: pipeliner}
			if e := p.Set(ctx, fmt.Sprintf(defaultHandleToIdResolver, handle), u.ID.String(), time.Hour*12).Err(); e != nil {
				return e
			}
			k := fmt.Sprintf(defaultUserKey, u.ID)
			if e := p.JSONSet(ctx, k, "$", u); e != nil {
				return e
			}
			return s.j.Expire(ctx, k, time.Hour*12).Err()
		})
		return u
	}
	return nil
}

func (s *UserStore) Get(ctx context.Context, id uuid.UUID, handle string) *user.User {
	if handle == "" && id == uuid.Nil {
		return nil
	}
	if id == uuid.Nil {
		if _id, e := s.j.Get(ctx, fmt.Sprintf(defaultHandleToIdResolver, handle)).Result(); e == nil && _id != "" && _id != uuid.Nil.String() {
			id, e = uuid.Parse(_id)
			if e == nil {
				return s.fallback(ctx, id, "")
			}
		}
		return s.fallback(ctx, uuid.Nil, handle)
	}
	//if !s.Exists(ctx, id) {
	//	return nil
	//}
	r := s.j.JSONGet(ctx, fmt.Sprintf(defaultUserKey, id), "$")
	if _u := rejson.Unmarshal[user.User](r); _u == nil {
		return s.fallback(ctx, id, "")
	} else if len(_u) > 0 {
		return &(_u)[0]
	}
	return nil
}

func userToMinimalUser(u *user.User) *user.MinimalUser {
	if u == nil {
		return nil
	}
	var topRole interface{} = nil
	if len(u.Roles) > 0 {
		topRole = u.Roles[0]
	}
	return &user.MinimalUser{
		ID:           u.ID.String(),
		DisplayName:  u.DisplayName,
		Handle:       u.Handle,
		Avatar:       u.Avatar,
		Organization: u.Organization,
		TopRole:      topRole,
		Rating:       u.Rating,
	}
}

func (s *UserStore) GetMinimal(ctx context.Context, id uuid.UUID) *user.MinimalUser {
	if id == uuid.Nil {
		return nil
	}
	//if !s.Exists(ctx, id) {
	//	return nil
	//}
	r := s.j.JSONGet(ctx, fmt.Sprintf(defaultUserKey, id), "$['id','displayName','handle','organization','avatar','roles','rating']")
	res := rejson.Unmarshal[interface{}](r)
	if res == nil || len(res) != 7 {
		return userToMinimalUser(s.fallback(ctx, id, ""))
	}
	var topRole interface{} = nil
	if _r, ok := res[5].([]interface{}); ok && len(_r) > 0 {
		topRole = _r[0]
	}
	return &user.MinimalUser{
		ID:           res[0].(string),
		DisplayName:  res[1].(string),
		Handle:       res[2].(string),
		Organization: res[3].(string),
		Avatar:       res[4].(string),
		TopRole:      topRole,
		Rating:       uint16(res[6].(float64)),
	}
}
