package limiter

import (
	"sync"
	"time"
)

type SlidingWindowCounter struct {
	currentWindow int64
	currentCount  int64
	previousCount int64
	mu            sync.Mutex
}

type SlidingWindowLimiter struct {
	counters   map[string]*SlidingWindowCounter
	mu         sync.RWMutex
	limit      int64
	windowSize time.Duration
}

func (l *SlidingWindowLimiter) GetCounter(ip string) *SlidingWindowCounter {
	l.mu.RLock()
	swc, exists := l.counters[ip]
	l.mu.RUnlock()

	if exists {
		return swc
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	swc, exists = l.counters[ip]

	if exists {
		return swc
	}
	swc = &SlidingWindowCounter{}
	l.counters[ip] = swc

	return swc

}
func NewSlidingWindowLimitor(limit int64, windowSize time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		counters:   map[string]*SlidingWindowCounter{},
		limit:      limit,
		windowSize: windowSize,
	}
}
func (l *SlidingWindowLimiter) Allow(ip string) bool {
	counter := l.GetCounter(ip)
	counter.mu.Lock()
	defer counter.mu.Unlock()

	windowNumber := time.Now().Unix() / int64(l.windowSize.Seconds())

	if windowNumber != counter.currentWindow {
		if windowNumber == counter.currentWindow+1 {
			counter.previousCount = counter.currentCount
		} else {
			counter.previousCount = 0
		}
		counter.currentCount = 0
		counter.currentWindow = windowNumber
	}
	windowSizeSeconds := int64(l.windowSize.Seconds())
	windowStartUnix := windowNumber * windowSizeSeconds
	windowStartTime := time.Unix(windowStartUnix, 0)
	timeElapsed := time.Since(windowStartTime)

	overlapWeight := float64(l.windowSize-timeElapsed) / float64(l.windowSize)
	estimatedCount := float64(counter.currentCount) + (float64(counter.previousCount) * overlapWeight)

	if estimatedCount < float64(l.limit) {
		counter.currentCount++
		return true
	}
	return false
}
