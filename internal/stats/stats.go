package stats

import (
	"sync"
	"sync/atomic"
	"time"
)

type SuffixStats struct {
	TotalRequests atomic.Int64
	TotalErrors   atomic.Int64
	recentCounts  []int64
	mu            sync.Mutex
}

type Collector struct {
	mu      sync.RWMutex
	suffixes map[string]*SuffixStats
	globalQPS atomic.Int64
	ticker    *time.Ticker
	done      chan struct{}
}

var Global = &Collector{
	suffixes: make(map[string]*SuffixStats),
	done:     make(chan struct{}),
}

func (c *Collector) Start() {
	c.ticker = time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-c.ticker.C:
				c.mu.RLock()
				var total int64
				for _, s := range c.suffixes {
					current := s.TotalRequests.Load()
					s.mu.Lock()
					s.recentCounts = append(s.recentCounts, current)
					if len(s.recentCounts) > 60 {
						s.recentCounts = s.recentCounts[1:]
					}
					s.mu.Unlock()
					if len(s.recentCounts) >= 2 {
						total += s.recentCounts[len(s.recentCounts)-1] - s.recentCounts[len(s.recentCounts)-2]
					}
				}
				c.globalQPS.Store(total)
				c.mu.RUnlock()
			case <-c.done:
				return
			}
		}
	}()
}

func (c *Collector) Stop() {
	if c.ticker != nil {
		c.ticker.Stop()
	}
	close(c.done)
}

func (c *Collector) Record(suffix string, isError bool) {
	s := c.getOrCreate(suffix)
	s.TotalRequests.Add(1)
	if isError {
		s.TotalErrors.Add(1)
	}
}

func (c *Collector) getOrCreate(suffix string) *SuffixStats {
	c.mu.RLock()
	s, ok := c.suffixes[suffix]
	c.mu.RUnlock()
	if ok {
		return s
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	s = &SuffixStats{}
	c.suffixes[suffix] = s
	return s
}

func (c *Collector) GetQPS() int64 {
	return c.globalQPS.Load()
}

type SuffixSummary struct {
	Suffix   string `json:"suffix"`
	Requests int64  `json:"requests"`
	Errors   int64  `json:"errors"`
}

func (c *Collector) GetSuffixSummaries() []SuffixSummary {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]SuffixSummary, 0, len(c.suffixes))
	for suffix, s := range c.suffixes {
		result = append(result, SuffixSummary{
			Suffix:   suffix,
			Requests: s.TotalRequests.Load(),
			Errors:   s.TotalErrors.Load(),
		})
	}
	return result
}
