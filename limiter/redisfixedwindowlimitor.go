package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redisfixedwindowlimitor struct {
	client     *redis.Client
	limit      int64
	windowsize time.Duration
}

func NewRedisFixedWindowLimitor(limit int64, windowsize time.Duration, client *redis.Client) *Redisfixedwindowlimitor {
	return &Redisfixedwindowlimitor{
		client:     client,
		limit:      limit,
		windowsize: windowsize,
	}
}

func (l *Redisfixedwindowlimitor) Allow(ip string) bool {
	ctx := context.Background()

	windowNumber := time.Now().Unix() / int64(l.windowsize.Seconds())

	key := fmt.Sprintf("ratelimit:%s:%d", ip, windowNumber)

	count, err := l.client.Incr(ctx, key).Result()

	if err != nil {
		return true
	}
	if count == 1 {
		l.client.Expire(ctx, key, l.windowsize+time.Second)
	}
	return count <= l.limit

}
