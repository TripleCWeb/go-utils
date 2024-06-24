package hash

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type HashManager struct {
	redis *redis.Client
}

func NewHashManager(client *redis.Client) *HashManager {
	return &HashManager{
		redis: client,
	}
}

func (hm *HashManager) Del(keys ...string) error {
	// 开启事务
	tx := hm.redis.TxPipeline()

	var deleteKeys []string
	for _, key := range keys {
		// 查找所有以key开头的key
		keys, err := hm.redis.Keys(context.Background(), fmt.Sprintf("%s*", key)).Result()
		if err != nil {
			return err
		}
		deleteKeys = append(deleteKeys, keys...)
	}

	res := tx.Del(context.Background(), deleteKeys...)

	// 执行事务
	_, err := tx.Exec(context.Background())
	if err != nil {
		return err
	}

	return res.Err()
}

func (hm *HashManager) HSet(key, field string, value interface{}) error {
	return hm.redis.HSet(context.Background(), key, field, value).Err()
}

func (hm *HashManager) HGet(key, field string) (string, error) {
	return hm.redis.HGet(context.Background(), key, field).Result()
}

func (hm *HashManager) HGetInt32(key, field string) (n32 int32, err error) {
	str, err := hm.redis.HGet(context.Background(), key, field).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	n, err := strconv.Atoi(str)
	n32 = int32(n)
	return
}

func (hm *HashManager) HGetString(key, field string) (str string, err error) {
	str, err = hm.redis.HGet(context.Background(), key, field).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return
}

func (hm *HashManager) HMSet(key string, fieldsAndValues ...interface{}) error {
	return hm.redis.HMSet(context.Background(), key, fieldsAndValues...).Err()
}

func (hm *HashManager) HMGet(key string, fields ...string) ([]interface{}, error) {
	return hm.redis.HMGet(context.Background(), key, fields...).Result()
}

func (hm *HashManager) HDel(key string, fields ...string) (int64, error) {
	return hm.redis.HDel(context.Background(), key, fields...).Result()
}

func (hm *HashManager) HGetAll(key string) (map[string]string, error) {
	return hm.redis.HGetAll(context.Background(), key).Result()
}

func (hm *HashManager) HKeys(key string) ([]string, error) {
	return hm.redis.HKeys(context.Background(), key).Result()
}

func (hm *HashManager) HVals(key string) ([]string, error) {
	return hm.redis.HVals(context.Background(), key).Result()
}

func (hm *HashManager) HLen(key string) (int64, error) {
	return hm.redis.HLen(context.Background(), key).Result()
}

func (hm *HashManager) HExists(key, field string) (bool, error) {
	return hm.redis.HExists(context.Background(), key, field).Result()
}

func (hm *HashManager) HIncrBy(key, field string, incr int64) (int64, error) {
	return hm.redis.HIncrBy(context.Background(), key, field, incr).Result()
}

func (hm *HashManager) HIncrByFloat(key, field string, incr float64) (float64, error) {
	return hm.redis.HIncrByFloat(context.Background(), key, field, incr).Result()
}
