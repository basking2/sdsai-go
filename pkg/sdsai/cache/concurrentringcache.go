package cache

import (
	"hash/crc32"
	"math"
	"sync"
)

// A cache that is comprised other caches in a ring.
//
// Each cache is protected with a Mutex that is locked for write
// operations.
//
// Keys are hashed into the ring and added to their respective
// caches.
//
type ConcurrentRingCache struct {
	SizeLimit int
	RingSize  int
	Caches    []*LIFOCache
	Locks     []sync.RWMutex

	// A function that hashes a string into one of the caches in the Caches array.
	// The ringSize is the length of the Cache and Locks arrays.
	KeyHash func(key string, ringSize int) int
}

func NewConcurrentRingCache(ringSize int) *ConcurrentRingCache {
	c := ConcurrentRingCache{}

	c.RingSize = ringSize
	c.Caches = make([]*LIFOCache, ringSize)
	c.Locks = make([]sync.RWMutex, ringSize)
	c.KeyHash = func(s string, ringSize int) int {
		h := int(crc32.ChecksumIEEE([]byte(s)))

		if h < 0 {
			h = -1
		}

		h = h % ringSize

		return h
	}

	for i := 0; i < ringSize; i++ {
		c.Caches[i] = NewLIFOCache()
		c.Locks[i] = sync.RWMutex{}
	}

	return &c
}

// Add an item and enforce the cache size limit.
func (c *ConcurrentRingCache) Put(key string, item interface{}) {
	h := c.KeyHash(key, c.RingSize)

	c.Locks[h].Lock()
	c.Caches[h].Put(key, item)
	c.Locks[h].Unlock()
}

func (c *ConcurrentRingCache) PutWithHandler(key string, item interface{}, evictionHandler func(string, interface{})) {
	h := c.KeyHash(key, c.RingSize)

	c.Locks[h].Lock()
	c.Caches[h].PutWithHandler(key, item, evictionHandler)
	c.Locks[h].Unlock()
}

func (c *ConcurrentRingCache) Get(key string) (interface{}, bool) {
	h := c.KeyHash(key, c.RingSize)

	c.Locks[h].RLock()
	defer c.Locks[h].RUnlock()
	return c.Caches[h].Get(key)
}

// Evict items from every sub-cache until they contain the ceiling of 1/N
// items where N is the size limit for this entire cache.
func (c *ConcurrentRingCache) EnforceDistributedSizeLimit() {
	limit := int(math.Ceil(float64(c.SizeLimit) / float64(c.RingSize)))

	for i := 0; i < c.RingSize; i++ {
		c.Locks[i].Lock()
		for c.Caches[i].Len() > limit {
			c.Caches[i].EvictNext()
		}
		c.Locks[i].Unlock()
	}
}

func (c *ConcurrentRingCache) Size() int {
	size := 0

	for i := 0; i < c.RingSize; i++ {
		c.Locks[i].RLock()
		size += c.Caches[i].Len()
		c.Locks[i].RUnlock()
	}

	return size
}
