package detector

import (
	"sync"
	"time"
)

type BlockList struct {
	blocked map[string]time.Time
	mu      sync.RWMutex
}

func (b *BlockList) Block(ip string, duration time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.blocked[ip] = time.Now().Add(duration)
}

func (b *BlockList) IsBlocked(ip string) bool {

	b.mu.RLock()
	defer b.mu.RUnlock()

	expiry, exists := b.blocked[ip]

	if !exists {
		return false
	}
	return time.Now().Before(expiry)
}
