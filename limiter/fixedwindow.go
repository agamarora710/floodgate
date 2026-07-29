package limiter

import (
	"sync"
	"time"
)

type FixedWindowCounter struct {
	limit       int64
	windowsize  time.Duration
	windowstart int64
	count       int64
	mu          sync.Mutex
}

type FixedWindowLimitor struct {
	counters   map[string]*FixedWindowCounter
	mu         sync.RWMutex
	limit      int64
	windowSize time.Duration
}

func (l *FixedWindowLimitor) GetCounter(ip string) *FixedWindowCounter {
	l.mu.RLock()
	fwc, exists := l.counters[ip]
	l.mu.RUnlock()

	if exists {
		return fwc
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	fwc, exists = l.counters[ip]
	if exists {
		return fwc
	}
	fwc = &FixedWindowCounter{
		limit:      l.limit,
		windowsize: l.windowSize,
		count:      0,
	}
	l.counters[ip] = fwc

	return fwc
}

func NewFixedWindowLimitor(limit int64, windowSize time.Duration) *FixedWindowLimitor {
	return &FixedWindowLimitor{
		counters:   map[string]*FixedWindowCounter{},
		limit:      limit,
		windowSize: windowSize,
	}
}
func (l *FixedWindowLimitor) Allow(ip string) bool {
	counter := l.GetCounter(ip)
	counter.mu.Lock()
	defer counter.mu.Unlock()

	currentWindow := time.Now().Unix() / int64(l.windowSize.Seconds()) //gives current window number we are in !

	if currentWindow != counter.windowstart {
		counter.windowstart = currentWindow
		counter.count = 0
	}

	if counter.count < l.limit {
		counter.count++
		return true
	}

	return false

}
