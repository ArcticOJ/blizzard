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

func args(cmd Command, key string, v ...interface{}) []interface{} {
	return append([]interface{}{string(cmd), key}, v...)
}

func (r *JsonResult) String() string {
	return string(r.data)
}

func (r *JsonResult) Raw() json.RawMessage {
	return r.data
}

func Unmarshal[T any](r *JsonResult) []T {
	if r == nil {
		return nil
	}
	var obj []T
	if json.Unmarshal(r.data, &obj) != nil {
		return nil
	}
	return obj
}

func (r *ReJSON) JSONGet(ctx context.Context, key string, paths ...interface{}) *JsonResult {
	s, e := r.Do(ctx, args(GET, key, paths...)...).Result()
	if e != nil || len(s.(string)) == 0 {
		return nil
	}
	return &JsonResult{
		data: []byte(s.(string)),
	}
}

func (r *ReJSON) JSONSet(ctx context.Context, key string, path string, data interface{}) error {
	b, e := json.Marshal(data)
	if e != nil {
		return e
	}
	return r.Do(ctx, args(SET, key, path, string(b))...).Err()
}
