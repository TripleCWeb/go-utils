package rcache

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/TripleCGame/apis/pkg/utils/json"
	"github.com/redis/go-redis/v9"
)

type RCache interface {
	Set(ctx context.Context, key string, in interface{}, out interface{}) error
	Get(ctx context.Context, key string, in interface{}, out interface{}) error
	Exist(ctx context.Context, key string, in interface{}) bool
	DelField(ctx context.Context, key string, ins ...interface{}) error
	DelKey(ctx context.Context, keys ...string) error
	FlushAll(ctx context.Context) error
	FlushDB(ctx context.Context) error
}

func RcacheWapper[InType, OutType any](rcache RCache, key string, f func(context.Context, InType) (OutType, error), useCache bool) func(context.Context, InType) (OutType, error) {
	return func(ctx context.Context, in InType) (out OutType, err error) {
		funcNames := strings.Split(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), ".")
		funcName := strings.Split(funcNames[len(funcNames)-1], "-")[0]
		inValue := json.Interface2String(in)
		rcacheKey := fmt.Sprintf("%s:%s:%s", key, funcName, inValue)

		// Attempt to retrieve the response from the cache
		if useCache {
			err = rcache.Get(context.Background(), rcacheKey, in, &out)
			if err != redis.Nil {
				return
			}
		}

		// If the response was not found in the cache, call the underlying RPC method
		out, err = f(ctx, in)
		if err != nil {
			return
		}

		// write it to the cache
		if useCache {
			err = rcache.Set(context.Background(), rcacheKey, in, out)
			if err != nil {
				return
			}
		}

		return out, nil
	}
}
