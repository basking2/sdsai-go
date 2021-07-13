package cache

import (
	"fmt"
	"testing"
)

func TestConcurrentRingCache(t *testing.T) {
	cache := NewConcurrentRingCache(1)

	cache.SizeLimit = 5
	evicted := 0

	for i := 0; i < 100; i++ {
		cache.PutWithHandler(fmt.Sprintf("key %d", i), i, func(key string, obj interface{}) {
			evicted += 1
		})

		cache.EnforceDistributedSizeLimit()
	}

	if evicted != 95 {
		t.Errorf("Expected 95 but found %d evicted.", evicted)
	}
}
