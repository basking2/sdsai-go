package cache

import (
	"hash/crc32"
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
	AgeLimit  int64
	Caches    []*LIFOCache
	Locks     []*sync.RWMutex

	// A function that hashes a string into one of the caches in the Caches array.
	// The ringSize is the length of the Cache and Locks arrays.
	KeyHash func(key string, ringSize int) int
}

// Create a new concurrent ring cache.
// ringSize is how many independent caches will be created.
// cacheSize is how large each individual cache may be.
// ageLimit is how old an item may be if it may be returned.
//          If an item is fetched that is older than the ageLimit,
//          it will not be returned and the cache it resides in will be
//          updated to expire all older items.
//          If this is less than 0, no limit is applied.
func NewConcurrentRingCache(ringSize int, cacheSize int, ageLimit int64) *ConcurrentRingCache {
	c := ConcurrentRingCache{}

	c.RingSize = ringSize
	c.SizeLimit = cacheSize
	c.AgeLimit = ageLimit
	c.Caches = make([]*LIFOCache, ringSize)
	c.Locks = make([]*sync.RWMutex, ringSize)
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
		c.Locks[i] = &sync.RWMutex{}
	}

	return &c
}

// Set the time function that each cache in the ring of caches will use.
func (c *ConcurrentRingCache) SetTimeFunction(timeFunction func() int64) {
	c.EachSubCache(func(c *LIFOCache) {
		c.TimeFunction = timeFunction
	})
}

// Lock each sub-cache and pass it to the handler function.
//
// Each cache is locked as writable first.
func (c *ConcurrentRingCache) EachSubCache(f func(*LIFOCache)) {
	for i := 0; i < c.RingSize; i++ {
		c.Locks[i].Lock()

		f(c.Caches[i])

		c.Locks[i].Unlock()
	}
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

// Get an item from the sub-cache that holds items for the given key.
//
// If the key is not found in the sub-cache, (nil, false) is returned.
//
// If the key is found in the sub-cache but it is expired, (item, false) is
// returned where the item is the the expired data. The sub-cache
// is also cleaned so that no item older than the AgeLimit remains.
//
// If the key is found in the sub-cache and it is not expired, (item, true)
// is returned where the item is the user's data.
func (c *ConcurrentRingCache) Get(key string) (interface{}, bool) {
	h := c.KeyHash(key, c.RingSize)

	c.Locks[h].RLock()
	defer c.Locks[h].RUnlock()
	item, addedAt, ok := c.Caches[h].Get(key)

	if !ok {
		return nil, false
	}

	// If there is an age limit...
	if c.AgeLimit >= 0 {

		// If the item is older than the age limit (using the cache's time function to get "now")...
		if c.Caches[h].TimeFunction()-addedAt > c.AgeLimit {

			// Clean up only this cache...
			c.Caches[h].EvictOlderThan(c.AgeLimit)

			// And return that we couldn't find the item.
			// NOTE: Even if expired, we do return the found item.
			//       This gives the user more options.
			return item, false
		}
	}

	return item, true
}

// Evict items from every sub-cache until they contain the ceiling of 1/N
// items where N is the size limit for this entire cache.
func (c *ConcurrentRingCache) EnforceSizeLimit() {
	limit := c.SizeLimit

	c.EachSubCache(func(c *LIFOCache) {
		for c.Len() > limit {
			c.EvictNext()
		}
	})
}

func (c *ConcurrentRingCache) EvictOrderThan(tm int64) {
	c.EachSubCache(func(c *LIFOCache) {
		c.EvictOlderThan(tm)
	})
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
