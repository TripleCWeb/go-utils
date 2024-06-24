package lock

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "host.docker.internal:6379",
		Password: "rdspwd11131456",
		DB:       9,
	})
	return client
}

func currentTime() string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	return formattedTime
}

func TestGetLock(t *testing.T) {
	l := NewLocker(NewRedis())
	key := "test_key"
	timeout := 10

	isLock := l.AcquireLock(key, timeout)
	if !isLock {
		t.Error("Failed to acquire lock")
	}
}

func TestGetLockErr(t *testing.T) {
	l := NewLocker(NewRedis())
	key := "test_key"
	timeout := 1

	err := l.ReleaseLock(key)
	if err != nil {
		t.Errorf("Error occurred while releasing lock: %v", err)
	}

	isLock := l.AcquireLock(key, timeout)
	if !isLock {
		t.Error("Failed to acquire lock")
	}

	isLock = l.AcquireLock(key, timeout)
	if isLock {
		t.Error("Success to acquire lock")
	}
}

func TestReleaseLock(t *testing.T) {
	l := NewLocker(NewRedis())
	key := "test_key"

	err := l.ReleaseLock(key)
	if err != nil {
		t.Errorf("Error occurred while releasing lock: %v", err)
	}
}

func TestLockOperations_MultipleSessions(t *testing.T) {
	// 创建多个会话
	sessionCount := 5
	var wg sync.WaitGroup
	wg.Add(sessionCount)

	// 使用互斥锁保护交互操作
	var mutex sync.Mutex
	successCount := 0
	timeout := 20

	key := "TestLockOperations_MultipleSessions"

	for i := 0; i < sessionCount; i++ {
		go func(sessionID int) {
			defer wg.Done()

			l := NewLocker(NewRedis())

			isLock := l.AcquireLock(key, timeout)
			if !isLock {
				t.Errorf("Failed to acquire lock (session %d)", sessionID)
				return
			}

			fmt.Printf("%s Session %d successfully acquired lock, sleep 1 seconds\n", currentTime(), sessionID)
			time.Sleep(1 * time.Second)

			// 使用互斥锁保护交互操作
			mutex.Lock()
			successCount++
			fmt.Printf("%s Session %d prepare released lock\n", currentTime(), sessionID)
			mutex.Unlock()

			// 释放锁
			err := l.ReleaseLock(key)
			if err != nil {
				t.Errorf("Error occurred while releasing lock (session %d): %v", sessionID, err)
				return
			}

			// 使用互斥锁保护交互操作
			mutex.Lock()
			fmt.Printf("%s Session %d successfully released lock\n", currentTime(), sessionID)
			mutex.Unlock()
		}(i)
	}

	wg.Wait()

	// 验证成功获取锁的会话数量
	if successCount != sessionCount {
		t.Errorf("%s Failed to acquire lock in %d out of %d sessions", currentTime(), sessionCount-successCount, sessionCount)
	}
}
