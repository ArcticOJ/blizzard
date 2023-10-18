package rejson

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

type (
	Client interface {
		redis.Cmdable
		Do(ctx context.Context, args ...interface{}) *redis.Cmd
	}
	ReJSON struct {
		Client
	}
	JsonResult struct {
		data []byte
	}
)

func args(cmd Command, key string, paths []interface{}, v ...interface{}) []interface{} {
	return append(append([]interface{}{cmd, key}, paths...), v...)
}

func (r JsonResult) String() string {
	return string(r.data)
}

func (r JsonResult) Raw() json.RawMessage {
	return r.data
}

func Unmarshal[T any](r JsonResult) *T {
	if len(r.data) == 0 {
		return nil
	}
	var obj T
	if json.Unmarshal(r.data, &obj) != nil {
		return nil
	}
	return &obj
}

func (r *ReJSON) JSONGet(ctx context.Context, key string, paths ...interface{}) JsonResult {
	s, e := r.Do(ctx, args(GET, key, paths)...).Result()
	if e != nil || len(s.(string)) == 0 {
		return JsonResult{}
	}
	return JsonResult{
		data: []byte(s.(string)),
	}
}

func (r *ReJSON) JSONSet(ctx context.Context, key string, path string, data interface{}) error {
	b, e := json.Marshal(data)
	if e != nil {
		return e
	}
	return r.Do(ctx, args(SET, key, []interface{}{path}, string(b))...).Err()
}

func (r *ReJSON) JSONMGet(ctx context.Context, path string, keys ...interface{}) (res []JsonResult) {
	s, e := r.Do(ctx, append([]interface{}{MGET}, append(keys, path)...)...).Result()
	if e != nil {
		return
	}
	_r, ok := s.([]interface{})
	if ok && len(_r) > 0 {
		for _, _s := range _r {
			if _s == nil {
				res = append(res, JsonResult{})
			} else {
				res = append(res, JsonResult{data: []byte(_s.(string))})
			}
		}
	}
	return
}

func (r *ReJSON) JTxPipelined(ctx context.Context, fn func(r *ReJSON) error) error {
	_, e := r.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		return fn(&ReJSON{Client: pipeliner})
	})
	return e
}

func (r *ReJSON) JTxPipeline(ctx context.Context) *ReJSON {
	return &ReJSON{Client: r.TxPipeline()}
}
