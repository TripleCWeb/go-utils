package rcache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
)

const ExpirationSeconds time.Duration = 300 * time.Second

type rcache struct {
	redis *redis.Client
	log   *log.Helper
}

func New(redis *redis.Client, log *log.Helper) RCache {
	return &rcache{
		redis: redis,
		log:   log,
	}
}

func (r *rcache) Set(ctx context.Context, key string, in interface{}, out interface{}) error {
	field := FieldSerializer{}.Serialize(in)
	value := ValueSerializer{}.Serialize(out)

	// 开启事务
	tx := r.redis.TxPipeline()

	// 执行HSet操作
	tx.HSet(ctx, key, field, value)

	// 执行Expire操作
	tx.Expire(ctx, key, ExpirationSeconds)

	// 执行事务
	_, err := tx.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
func (r *rcache) Get(ctx context.Context, key string, in interface{}, out interface{}) error {
	field := FieldSerializer{}.Serialize(in)
	cacheStr, err := r.redis.HGet(ctx, key, field).Result()
	if err != nil {
		return err
	}
	return ValueSerializer{}.Deserialize(cacheStr, out)
}
func (r *rcache) Exist(ctx context.Context, key string, in interface{}) bool {
	field := FieldSerializer{}.Serialize(in)
	isExisted, _ := r.redis.HExists(ctx, key, field).Result()
	return isExisted
}
func (r *rcache) DelField(ctx context.Context, key string, ins ...interface{}) error {
	var fields []string
	for _, in := range ins {
		fields = append(fields, FieldSerializer{}.Serialize(in))
	}
	res := r.redis.HDel(ctx, key, fields...)
	return res.Err()
}
func (r *rcache) DelKey(ctx context.Context, keys ...string) error {
	// 开启事务
	tx := r.redis.TxPipeline()

	var deleteKeys []string
	for _, key := range keys {
		// 查找所有以key开头的key
		keys, err := r.redis.Keys(ctx, fmt.Sprintf("%s:*", key)).Result()
		if err != nil {
			return err
		}
		deleteKeys = append(deleteKeys, keys...)
	}

	res := tx.Del(ctx, deleteKeys...)

	// 执行事务
	_, err := tx.Exec(ctx)
	if err != nil {
		return err
	}

	return res.Err()
}

func (r *rcache) FlushAll(ctx context.Context) error {
	res := r.redis.FlushAll(ctx)
	return res.Err()
}

func (r *rcache) FlushDB(ctx context.Context) error {
	res := r.redis.FlushDB(ctx)
	return res.Err()
}
