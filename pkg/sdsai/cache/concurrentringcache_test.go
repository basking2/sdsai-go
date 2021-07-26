package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestConcurrentRingCache(t *testing.T) {
	cache := NewConcurrentRingCache(1, 5, -1)

	evicted := 0

	for i := 0; i < 100; i++ {
		cache.PutWithHandler(fmt.Sprintf("key %d", i), i, func(key string, obj interface{}) {
			evicted += 1
		})

		cache.EnforceSizeLimit()
	}

	if evicted != 95 {
		t.Errorf("Expected 95 but found %d evicted.", evicted)
	}
}

func TestConcurrentRingCacheConcurrent(t *testing.T) {
	cache := NewConcurrentRingCache(10, 5, -1)

	wg := sync.WaitGroup{}
	wg.Add(1000)

	for i := 0; i < 1000; i++ {

		go func(i int) {
			cache.Put(fmt.Sprintf("key %d", i), i)

			cache.EnforceSizeLimit()

			wg.Done()
		}(i)
	}

	wg.Wait()
}

func TestConcurrentRingCacheTimeEvict(t *testing.T) {
	cache := NewConcurrentRingCache(10, 5, 100)
	called := false
	cache.SetTimeFunction(func() int64 {
		return time.Now().UnixNano()
	})
	cache.PutWithHandler("hi", "hellooooo", func(string, interface{}) { called = true })

	cache.EvictOrderThan(time.Now().UnixNano() + 1000)

	if !called {
		t.Errorf("Did not call eviction callback.")
	}

}
