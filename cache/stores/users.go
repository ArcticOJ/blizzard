package stores

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/ArcticOJ/blizzard/v0/cache"
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/models/user"
	"github.com/ArcticOJ/blizzard/v0/logger"
	"github.com/ArcticOJ/blizzard/v0/rejson"
	"github.com/ArcticOJ/blizzard/v0/utils/numeric"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/tmthrgd/go-hex"
	"github.com/uptrace/bun"
	"golang.org/x/sync/singleflight"
	"strings"
	"time"
)

var Users *userStore

type userStore struct {
	j *rejson.ReJSON
	s singleflight.Group
}

const (
	defaultUserKey            = "blizzard::user[%s]"
	defaultHandleToIdResolver = "blizzard::user_id[%s]"
	defaultUserListKey        = "blizzard::user_list"
	defaultUserTtl            = time.Hour * 48

	DefaultUserPageSize = 25
)

func init() {
	Users = &userStore{j: &rejson.ReJSON{Client: cache.CreateClient(cache.User, "users")}}
	Users.populateUserList(context.Background())
}

func (s *userStore) populateUserList(c context.Context) {
	_, e, _ := s.s.Do("populate", func() (interface{}, error) {
		if Users.j.Exists(c, defaultUserListKey).Val() == 0 {
			var (
				ids     []string
				ratings []uint16
			)
			// TODO: use app's context for this
			if _e := db.Database.NewSelect().Model((*user.User)(nil)).Column("id", "rating").Scan(context.Background(), &ids, &ratings); _e != nil {
				return nil, _e
			}
			var m []redis.Z
			for i := range ids {
				m = append(m, redis.Z{
					Score:  float64(numeric.CompressUint16(ratings[i], 0)),
					Member: ids[i],
				})
			}
			if _e := Users.j.ZAdd(context.Background(), defaultUserListKey, m...).Err(); _e != nil {
				return nil, _e
			}
		}
		return nil, nil
	})
	logger.Panic(e, "failed to populate user cache")
}

func (s *userStore) loadOne(ctx context.Context, id uuid.UUID, handle string) (u *user.User) {
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
	if len(u.Roles) > 0 {
		u.TopRole = &u.Roles[0]
	}
	return
}

func (s *userStore) loadMulti(ctx context.Context, users []user.User) (u []user.User) {
	if db.Database.NewSelect().
		// users are not selected by order of ids but by order in the database, so I have to use this hack to ensure stable order
		// might not be performance-wise but at least it does the trick lol
		With("inp", db.Database.NewValues(&users).Column("id").WithOrder()).
		Model(&u).ExcludeColumn("password", "api_key").
		Table("inp").
		Where("\"user\".id = inp.id").
		OrderExpr("inp._order").
		Relation("Connections", func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.Where("show_in_profile = true")
		}).
		Relation("Roles", func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.Order("priority ASC").Column("name", "icon", "color")
		}).Scan(ctx) != nil {
		return nil
	}
	for i := range u {
		if strings.TrimSpace(u[i].Email) != "" {
			h := md5.Sum([]byte(u[i].Email))
			u[i].Avatar = hex.EncodeToString(h[:])
		}
		if len(u[i].Roles) > 0 {
			u[i].TopRole = &u[i].Roles[0]
		}
	}
	return
}

func (s *userStore) UserExists(ctx context.Context, id uuid.UUID) bool {
	s.populateUserList(ctx)
	return s.j.ZScore(ctx, defaultUserListKey, id.String()).Err() == nil
}

func (s *userStore) fallbackOne(ctx context.Context, id uuid.UUID, handle string) *user.User {
	u := s.loadOne(ctx, id, handle)
	if u != nil && u.ID != uuid.Nil {
		handle = u.Handle
		s.j.JTxPipelined(ctx, func(r *rejson.ReJSON) error {
			if e := r.Set(ctx, fmt.Sprintf(defaultHandleToIdResolver, handle), u.ID.String(), defaultUserTtl).Err(); e != nil {
				return e
			}
			k := fmt.Sprintf(defaultUserKey, u.ID)
			if e := r.JSONSet(ctx, k, "$", u); e != nil {
				return e
			}
			return r.Expire(ctx, k, defaultUserTtl).Err()
		})
		return u
	}
	return nil
}

func (s *userStore) fallbackMulti(ctx context.Context, users []user.User) (u []user.User) {
	u = s.loadMulti(ctx, users)
	s.j.JTxPipelined(ctx, func(r *rejson.ReJSON) error {
		for _, _u := range u {
			if _u.ID != uuid.Nil {
				k := fmt.Sprintf(defaultUserKey, _u.ID)
				if r.Set(ctx, fmt.Sprintf(defaultHandleToIdResolver, _u.Handle), _u.ID.String(), defaultUserTtl).Err() == nil &&
					r.JSONSet(ctx, k, "$", _u) == nil {
					r.Expire(ctx, k, defaultUserTtl)
				}
			}
		}
		return nil
	})
	return
}

func (s *userStore) Get(ctx context.Context, id uuid.UUID, handle string) *user.User {
	if handle == "" && id == uuid.Nil {
		return nil
	}
	if id == uuid.Nil {
		if _id, e := s.j.Get(ctx, fmt.Sprintf(defaultHandleToIdResolver, handle)).Result(); e == nil && _id != "" && _id != uuid.Nil.String() {
			id, e = uuid.Parse(_id)
			if e != nil {
				return s.fallbackOne(ctx, id, "")
			}
			goto getFromCache
		}
		return s.fallbackOne(ctx, uuid.Nil, handle)
	}
	if !s.UserExists(ctx, id) {
		return nil
	}
getFromCache:
	r := s.j.JSONGet(ctx, fmt.Sprintf(defaultUserKey, id))
	_u := rejson.Unmarshal[user.User](r)
	if _u == nil {
		return s.fallbackOne(ctx, id, "")
	}
	return _u
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

func (s *userStore) GetMinimal(ctx context.Context, id uuid.UUID) *user.MinimalUser {
	if id == uuid.Nil {
		return nil
	}
	if !s.UserExists(ctx, id) {
		return nil
	}
	r := s.j.JSONGet(ctx, fmt.Sprintf(defaultUserKey, id), "$..['id','displayName','handle','organization','avatar','topRole','rating']")
	res := rejson.Unmarshal[[]interface{}](r)
	if res == nil || len(*res) != 7 {
		return userToMinimalUser(s.fallbackOne(ctx, id, ""))
	}
	_res := *res
	return &user.MinimalUser{
		ID:           _res[0].(string),
		DisplayName:  _res[1].(string),
		Handle:       _res[2].(string),
		Organization: _res[3].(string),
		Avatar:       _res[4].(string),
		TopRole:      _res[5],
		Rating:       uint16(_res[6].(float64)),
	}
}

func (s *userStore) GetPage(ctx context.Context, page uint16, rev bool) (res []user.MinimalUser) {
	order := "desc"
	if rev {
		order = "asc"
	}
	r, e, _ := s.s.Do(fmt.Sprintf("page-%d-%s", page, order), func() (interface{}, error) {
		s.populateUserList(ctx)
		u := s.j.ZRangeArgs(context.Background(), redis.ZRangeArgs{
			Key:     defaultUserListKey,
			Start:   "-inf",
			Stop:    "+inf",
			ByScore: true,
			// the leaderboard should be descending by default, so rev should make it ascending instead of descending
			Rev:    !rev,
			Offset: int64(page * DefaultUserPageSize),
			Count:  DefaultUserPageSize,
		}).Val()
		if len(u) == 0 {
			return nil, errors.New("no users")
		}
		var toGet []interface{}
		for _, z := range u {
			toGet = append(toGet, fmt.Sprintf(defaultUserKey, z))
		}
		if r := s.j.JSONMGet(ctx, "$..['id','displayName','handle','organization','avatar','topRole','rating']", toGet...); len(r) > 0 {
			res = make([]user.MinimalUser, len(r))
			var (
				toLoad  []user.User
				indices []int
			)
			for i := range r {
				_u := rejson.Unmarshal[[]interface{}](r[i])
				if _u == nil || len(*_u) != 7 {
					toLoad = append(toLoad, user.User{
						ID: uuid.MustParse(u[i]),
					})
					indices = append(indices, i)
					continue
				}
				usr := *_u
				res[i] = user.MinimalUser{
					ID:           usr[0].(string),
					DisplayName:  usr[1].(string),
					Handle:       usr[2].(string),
					Organization: usr[3].(string),
					Avatar:       usr[4].(string),
					TopRole:      usr[5],
					Rating:       uint16(usr[6].(float64)),
				}
			}
			if len(toLoad) > 0 {
				ul := s.fallbackMulti(ctx, toLoad)
				for i, c := range ul {
					if _u := userToMinimalUser(&c); _u != nil {
						res[indices[i]] = *_u
					}
				}
			}
			return res, nil
		}
		return nil, nil
	})
	if e != nil {
		return
	}
	if _r, ok := r.([]user.MinimalUser); ok {
		r = _r
	}
	return
}
