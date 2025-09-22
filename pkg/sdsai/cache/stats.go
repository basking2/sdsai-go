package cache

import "sync/atomic"

type CacheStats struct {
	hit   int32
	miss  int32
	evict int32
}

func (c *CacheStats) Hit() {
	atomic.AddInt32(&c.hit, 1)
}

func (c *CacheStats) Miss() {
	atomic.AddInt32(&c.miss, 1)
}

func (c *CacheStats) Evict() {
	atomic.AddInt32(&c.evict, 1)
}

func (c *CacheStats) Reset() {
	atomic.StoreInt32(&c.hit, 0)
	atomic.StoreInt32(&c.miss, 0)
	atomic.StoreInt32(&c.evict, 0)
}

// Return (hit, miss, evict) counts.
func (c *CacheStats) GetStats() (int32, int32, int32) {
	return atomic.LoadInt32(&c.hit), atomic.LoadInt32(&c.miss), atomic.LoadInt32(&c.evict)
}
