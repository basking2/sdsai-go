package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestConcurrentRingCache(t *testing.T) {
	cache, err := NewConcurrentRingCache(1, 5, -1)

	if err != nil {
		t.Fatal(err)
	}

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
	cache, err := NewConcurrentRingCache(10, 5, -1)
	if err != nil {
		t.Fatal(err)
	}

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
	cache, err := NewConcurrentRingCache(10, 5, 1)
	if err != nil {
		t.Fatal(err)
	}

	called := false
	cache.SetTimeFunction(func() int64 {
		return time.Now().UnixNano()
	})
	cache.PutWithHandler("hi", "hellooooo", func(string, interface{}) { called = true })

	d, _ := time.ParseDuration("10ms")
	time.Sleep(d)

	if item, ok := cache.Get("hi"); ok {
		t.Error("Item should not be present.")
	} else {
		switch s := item.(type) {
		case string:
			if s != "hellooooo" {
				t.Error("s is not hellooooos")
			}
		default:
			t.Error("Item is not a string.")
		}
	}

	if !called {
		t.Errorf("Did not call eviction callback.")
	}

}
