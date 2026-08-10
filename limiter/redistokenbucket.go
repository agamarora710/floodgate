package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const luatokencode = `
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local refillRate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

local tokens = tonumber(redis.call("HGET",key,"tokens"))
local lastRefill = tonumber(redis.call("HGET",key,"lastRefill"))

if tokens == nil then
    tokens = capacity
    lastRefill = now
end

local elapsed = now - lastRefill
local refillAmount = elapsed*refillRate

tokens = math.min(capacity,tokens+refillAmount)

local allowed = 0

if tokens>=1 then
    tokens = tokens - 1
    allowed = 1
end

redis.call("HSET",key,"tokens",tokens,"lastRefill",now)

return allowed
`

type RedisTokenBucketLimitor struct {
	client     *redis.Client
	capacity   float64
	refillRate float64
}

func NewRedisTokenBucketLimitor(capacity float64, refillRate float64, client *redis.Client) *RedisTokenBucketLimitor {
	return &RedisTokenBucketLimitor{
		capacity:   capacity,
		refillRate: refillRate,
		client:     client,
	}
}

func (l *RedisTokenBucketLimitor) Allow(ip string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("Bucket:%s", ip)
	now := time.Now().Unix()

	script := redis.NewScript(luatokencode)

	result, err := script.Run(
		ctx,
		l.client,
		[]string{key},
		l.capacity,
		l.refillRate,
		now,
	).Result()

	if err != nil {
		return true
	}
	return result.(int64) == 1
}
