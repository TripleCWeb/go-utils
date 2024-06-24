package hash

import (
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
	"launchpad.net/gocheck"
)

func Test(t *testing.T) {
	gocheck.TestingT(t)
}

type S struct{}

var _ = gocheck.Suite(&S{})

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "host.docker.internal:6379",
	Password: "rdspwd11131456",
	DB:       int(4),
})

// 设置哈希 + 获取哈希字段的值
func (s *S) TestSetGet(c *gocheck.C) {
	hashManager := NewHashManager(redisClient)

	err := hashManager.HSet("myhash", "field1", "value1")
	c.Assert(err, gocheck.IsNil)

	value, err := hashManager.HGet("myhash", "field1")
	c.Assert(err, gocheck.IsNil)
	fmt.Println("Field1 value:", value)

	c.Assert(value, gocheck.Equals, "value1")
}

// 同时设置多个哈希字段的值 + 同时获取多个哈希字段的值
func (s *S) TestSetMulti(c *gocheck.C) {
	hashManager := NewHashManager(redisClient)

	err := hashManager.HMSet("myhash", "field2", "value2", "field3", "value3")
	c.Assert(err, gocheck.IsNil)

	values, err := hashManager.HMGet("myhash", "field1", "field2", "field3")
	c.Assert(err, gocheck.IsNil)
	fmt.Println("Field values:", values)

	c.Assert(values[0], gocheck.Equals, "value1")
	c.Assert(values[1], gocheck.Equals, "value2")
	c.Assert(values[2], gocheck.Equals, "value3")
}

// 删除哈希字段
func (s *S) TestDelete(c *gocheck.C) {
	hashManager := NewHashManager(redisClient)

	err := hashManager.HMSet("myhash", "field2", "value2", "field3", "value3")
	c.Assert(err, gocheck.IsNil)

	numDeleted, err := hashManager.HDel("myhash", "field1", "field2")
	c.Assert(err, gocheck.IsNil)
	fmt.Println("Number of deleted fields:", numDeleted)

	values, err := hashManager.HMGet("myhash", "field1", "field2", "field3")
	c.Assert(err, gocheck.IsNil)
	fmt.Println("Field values:", values)

	c.Assert(values[0], gocheck.Equals, nil)
	c.Assert(values[1], gocheck.Equals, nil)
	c.Assert(values[2], gocheck.Equals, "value3")
}

// 获取所有哈希字段和值
func (s *S) TestGetAll(c *gocheck.C) {
	hashManager := NewHashManager(redisClient)

	err := hashManager.HMSet("myhash", "field1", "value1", "field2", "value2", "field3", "value3")
	c.Assert(err, gocheck.IsNil)

	allFields, err := hashManager.HGetAll("myhash")
	c.Assert(err, gocheck.IsNil)
	fmt.Println("All fields:", allFields) // 输出: All fields: map[field3:value3]

	c.Assert(len(allFields), gocheck.Equals, 3)
	for k, v := range allFields {
		if k == "field1" {
			c.Assert(v, gocheck.Equals, "value1")
		}
		if k == "field2" {
			c.Assert(v, gocheck.Equals, "value2")
		}
		if k == "field3" {
			c.Assert(v, gocheck.Equals, "value3")
		}
	}
}

// 获取哈希字段数量
func (s *S) TestLen(c *gocheck.C) {
	hashManager := NewHashManager(redisClient)

	err := hashManager.HMSet("myhash", "field1", "value1", "field2", "value2", "field3", "value3")
	c.Assert(err, gocheck.IsNil)

	numFields, err := hashManager.HLen("myhash")
	c.Assert(err, gocheck.IsNil)
	c.Assert(numFields, gocheck.Equals, int64(3))
}

// 检查哈希字段是否存在
func (s *S) TestHExist(c *gocheck.C) {
	hashManager := NewHashManager(redisClient)

	err := hashManager.HMSet("myhash", "field2", "value2", "field3", "value3")
	c.Assert(err, gocheck.IsNil)

	// 检查哈希字段是否存在
	exists, err := hashManager.HExists("myhash", "field3")
	c.Assert(err, gocheck.IsNil)
	c.Assert(exists, gocheck.Equals, true)

	numDeleted, err := hashManager.HDel("myhash", "field3")
	c.Assert(err, gocheck.IsNil)
	fmt.Println("Number of deleted fields:", numDeleted)

	// 检查哈希字段是否存在
	exists, err = hashManager.HExists("myhash", "field3")
	c.Assert(err, gocheck.IsNil)
	c.Assert(exists, gocheck.Equals, false)

}

// 哈希字段值增加整数
func (s *S) TestHIncrBy(c *gocheck.C) {
	hashManager := NewHashManager(redisClient)
	key := "myhashHIncrBy"
	hashManager.Del(key)

	// 哈希字段值增加整数
	newValue, err := hashManager.HIncrBy(key, "field3", 5)
	c.Assert(err, gocheck.IsNil)
	fmt.Println("New field3 value:", newValue) // 输出: New field3 value: 8

	value, err := hashManager.HGet(key, "field3")
	c.Assert(err, gocheck.IsNil)
	c.Assert(value, gocheck.Equals, "5")

	// 哈希字段值增加整数
	newValue, err = hashManager.HIncrBy(key, "field3", 10)
	c.Assert(err, gocheck.IsNil)
	fmt.Println("New field3 value:", newValue) // 输出: New field3 value: 8

	value, err = hashManager.HGet(key, "field3")
	c.Assert(err, gocheck.IsNil)
	c.Assert(value, gocheck.Equals, "15")
}

// 哈希字段值增加浮点数
func (s *S) TestHIncrByFloat(c *gocheck.C) {
	hashManager := NewHashManager(redisClient)
	key := "myhashHIncrByFloat"
	hashManager.Del(key)

	// 哈希字段值增加整数
	newValue, err := hashManager.HIncrByFloat(key, "field3", 5.1)
	c.Assert(err, gocheck.IsNil)
	fmt.Println("New field3 value:", newValue) // 输出: New field3 value: 8

	value, err := hashManager.HGet(key, "field3")
	c.Assert(err, gocheck.IsNil)
	c.Assert(value, gocheck.Equals, "5.1")

	// 哈希字段值增加整数
	newValue, err = hashManager.HIncrByFloat(key, "field3", 10.2)
	c.Assert(err, gocheck.IsNil)
	fmt.Println("New field3 value:", newValue) // 输出: New field3 value: 8

	value, err = hashManager.HGet(key, "field3")
	c.Assert(err, gocheck.IsNil)
	c.Assert(value, gocheck.Equals, "15.3")
}
