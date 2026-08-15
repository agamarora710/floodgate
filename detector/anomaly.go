package detector

import (
	"sync"
	"time"
)

type IPstats struct {
	recentCount   int64
	baselineCount int64
	windowStart   int64
	mu            sync.Mutex
}

type Detector struct {
	stats      map[string]*IPstats
	mu         sync.RWMutex
	windowSize time.Duration
	threshold  float64
}

func (d *Detector) GetStats(ip string) *IPstats {
	d.mu.RLock()
	detect, exists := d.stats[ip]
	d.mu.RUnlock()

	if exists {
		return detect
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if exists {
		return detect
	}
	detect = &IPstats{}

	d.stats[ip] = detect

	return detect
}

func (d *Detector) Record(ip string) {
	stats := d.GetStats(ip)
	stats.mu.Lock()
	defer stats.mu.Unlock()

	windowNumber := time.Now().Unix() / int64(d.windowSize.Seconds())

	if windowNumber != stats.windowStart {
		stats.baselineCount = stats.recentCount
		stats.recentCount = 0
		stats.windowStart = windowNumber
	}
	stats.recentCount++
}

func (d *Detector) IsAnamalous(ip string) bool {
	stats := d.GetStats(ip)
	stats.mu.Lock()
	defer stats.mu.Unlock()

	if stats.baselineCount == 0 {
		return false
	}
	ratio := float64(stats.recentCount) / float64(stats.baselineCount)

	if ratio > d.threshold {
		return true
	}
	return false
}
func NewDetector(windowSize time.Duration, threshold float64) *Detector {
	return &Detector{
		stats:      map[string]*IPstats{}
		windowSize: windowSize,
		threshold:  threshold,
	}
}
