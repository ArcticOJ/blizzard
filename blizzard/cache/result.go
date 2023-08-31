package cache

import (
	"blizzard/blizzard/db/models/contest"
	"blizzard/blizzard/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/rueidis"
	"strconv"
	"strings"
	"sync"
)

type ResultCache struct {
	rueidis.Client
	l  sync.RWMutex
	tl sync.RWMutex
}

var Result *ResultCache

const defaultResultKey = "blizzard::case_results[%d]"

const defaultTagKey = "blizzard::pending_submission[%d]"

const defaultTtl = 15

// 30 minutes
const defaultInitialTtl = 30 * 60

func init() {
	Result = &ResultCache{Client: create(1, "results")}
}

func (c *ResultCache) GetTag(id uint32) (uint64, error) {
	c.tl.RLock()
	defer c.tl.RUnlock()
	return c.Do(context.Background(), c.B().Get().Key(c.createTagKey(id)).Build()).AsUint64()
}

func (c *ResultCache) SetTag(id uint32, tag uint64) error {
	c.tl.Lock()
	defer c.tl.Unlock()
	key := c.createTagKey(id)
	res := c.DoMulti(
		context.Background(),
		c.B().Set().Key(key).Value(strconv.FormatUint(tag, 10)).Build(),
		c.B().Expire().Key(key).Seconds(defaultInitialTtl).Build())
	if e := res[0].Error(); e != nil {
		return e
	}
	return res[1].Error()
}

func (c *ResultCache) DeleteTag(id uint32) error {
	c.tl.Lock()
	defer c.tl.Unlock()
	return c.Do(context.Background(), c.B().Del().Key(c.createTagKey(id)).Build()).Error()
}

func (c *ResultCache) createKey(id uint32) string {
	return fmt.Sprintf(defaultResultKey, id)
}

func (c *ResultCache) createTagKey(id uint32) string {
	return fmt.Sprintf(defaultTagKey, id)
}

func (c *ResultCache) Create(id uint32, count uint16) error {
	if c.IsPending(id) {
		return errors.New("submission is being judged")
	}
	key := c.createKey(id)
	res := c.DoMulti(
		context.Background(),
		c.B().Del().Key(key).Build(),
		c.B().Rpush().Key(key).Element(utils.ArrayFill("", int(count))...).Build(),
		c.B().Expire().Key(key).Seconds(defaultInitialTtl).Build(),
	)
	if e := res[1].Error(); e != nil {
		return e
	}
	return res[2].Error()
}

func (c *ResultCache) IsPending(id uint32) bool {
	c.l.RLock()
	defer c.l.RUnlock()
	exists, e := c.Do(context.Background(), c.B().Exists().Key(c.createKey(id)).Build()).AsBool()
	return exists && e == nil
}

func (c *ResultCache) Store(id uint32, caseId uint16, r contest.CaseResult, ttl int) error {
	c.l.Lock()
	defer c.l.Unlock()
	buf, e := json.Marshal(r)
	if e != nil {
		return e
	}
	key := c.createKey(id)
	if ttl == 0 {
		ttl = defaultTtl
	}
	res := c.DoMulti(
		context.Background(),
		c.B().Lset().Key(key).Index(int64(caseId)).Element(string(buf)).Build(),
		c.B().Expire().Key(key).Seconds(int64(ttl)).Build(),
	)
	if e = res[0].Error(); e != nil {
		return e
	}
	return res[1].Error()
}

func (c *ResultCache) Get(id uint32) (string, error) {
	c.l.RLock()
	defer c.l.RUnlock()
	res := c.Do(context.Background(), c.B().Lrange().Key(c.createKey(id)).Start(0).Stop(-1).Build())
	if e := res.Error(); e != nil {
		return "[]", e
	}
	arr, e := res.AsStrSlice()
	if e != nil {
		return "[]", e
	}
	return fmt.Sprintf("[%s]", strings.Join(arr, ",")), nil
}

func (c *ResultCache) Clean(id uint32) {
	c.l.Lock()
	defer c.l.Unlock()
	c.Do(context.Background(), c.B().Del().Key(c.createKey(id)).Build())
}
