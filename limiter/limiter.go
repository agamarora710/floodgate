package limiter

import (
	"math"
	"sync"
	"time"
)

type Bucket struct {
	capacity   float64
	refillrate float64 //tokens per sec
	lastrefill time.Time
	tokens     float64
	mu         sync.Mutex
}

type TokenBucketLimtitor struct {
	buckets    map[string]*Bucket
	capacity   float64
	refillrate float64
	mu         sync.RWMutex
}

type Limiter interface {
	Allow(ip string) bool
}

func (l *TokenBucketLimtitor) GetBucket(ip string) *Bucket {
	l.mu.RLock()
	bucket, exists := l.buckets[ip]
	l.mu.RUnlock()

	if exists {
		return bucket
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	bucket, exists = l.buckets[ip]

	if exists {
		return bucket
	}

	bucket = &Bucket{
		capacity:   l.capacity,
		refillrate: l.refillrate,
		lastrefill: time.Now(),
		tokens:     l.capacity,
	}
	l.buckets[ip] = bucket

	return bucket

}

func (l *TokenBucketLimtitor) Allow(ip string) bool {

	bucket := l.GetBucket(ip)
	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	elapsedTime := time.Since(bucket.lastrefill)

	newTokens := elapsedTime.Seconds() * bucket.refillrate

	bucket.tokens = bucket.tokens + newTokens

	bucket.tokens = math.Min(bucket.capacity, bucket.tokens)

	bucket.lastrefill = time.Now()

	if bucket.tokens >= 1 {
		bucket.tokens--
		return true
	}
	return false

}

func NewTokenBucketLimitor(capacity, refillrate float64) *TokenBucketLimtitor {
	return &TokenBucketLimtitor{
		buckets:    map[string]*Bucket{},
		capacity:   capacity,
		refillrate: refillrate,
	}
}
