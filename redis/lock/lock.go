package lock

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Locker struct {
	client redis.UniversalClient
}

func NewLocker(client redis.UniversalClient) *Locker {
	return &Locker{client: client}
}

func (p *Locker) AcquireLock(lockKey string, waitTimeoutSecond time.Duration, lockTimeout time.Duration) bool {
	startTime := time.Now()
	endTime := startTime.Add(waitTimeoutSecond)

	for time.Now().Before(endTime) {
		ok, err := p.client.SetNX(context.Background(), lockKey, "locked", lockTimeout).Result()
		if err != nil {
			log.Println("Error acquiring lock:", err)
			return false
		}

		if ok {
			return true
		}

		time.Sleep(10 * time.Millisecond)
	}

	return false
}

func (p *Locker) ReleaseLock(lockKey string) error {
	_, err := p.client.Del(context.Background(), lockKey).Result()
	if err != nil {
		log.Println("Error releasing lock:", err)
	}
	return err
}
