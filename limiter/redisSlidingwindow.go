package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const luaslidingwindow = `
local key = KEYS[1]

local limit = tonumber(ARGV[1])
local windowSize =  tonumber(ARGV[2])
local now =  tonumber(ARGV[3])

local currentWindow = tonumber(redis.call("HGET",key,"currentWindow"))
local previousCount = tonumber(redis.call("HGET",key,"previousCount"))
local currentCount = tonumber(redis.call("HGET",key,"currentCount"))

if currentWindow == nil then
	currentWindow = 0
	previousCount = 0
	currentCount = 0
end

local windowNumber = math.floor(now/windowSize)

if currentWindow ~= windowNumber then
	if windowNumber == currentWindow + 1 then
		previousCount = currentCount
	else 
		previousCount = 0
	end
	currentCount = 0
	currentWindow = windowNumber
end

redis.call("EXPIRE", key, windowSize * 2)

local windowStart = currentWindow * windowSize
local elapsedIn = now - windowStart
local overlapWeight = (windowSize - elapsedIn)/windowSize
local estimatedCount = currentCount + (previousCount *overlapWeight)

local allowed = 0
if estimatedCount < limit then
	currentCount = currentCount + 1
	allowed = 1
end

redis.call("HSET", key, "currentWindow", currentWindow, "currentCount", currentCount, "previousCount", previousCount)

return allowed
`

type RedisSlidingWindowLimitor struct {
	client     *redis.Client
	limit      int64
	windowSize time.Duration
}

func NewRedisSlidingWindow(client *redis.Client, limit int64, windowSize time.Duration) *RedisSlidingWindowLimitor {
	return &RedisSlidingWindowLimitor{
		client:     client,
		limit:      limit,
		windowSize: windowSize,
	}
}

func (l *RedisSlidingWindowLimitor) Allow(ip string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("sliding:%s", ip)
	now := time.Now().UnixMilli()

	script := redis.NewScript(luaslidingwindow)

	result, err := script.Run(
		ctx,
		l.client,
		[]string{key},
		l.limit,
		l.windowSize.Milliseconds(),
		now,
	).Result()

	if err != nil {
		return true
	}
	return result.(int64) == 1
}
