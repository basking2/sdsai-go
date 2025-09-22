package integration

import (
	"testing"
	"time"

	"github.com/basking2/sdsai-go/pkg/sdsai/cache"
)

func TestExpiration(t *testing.T) {
	c, err := cache.NewConcurrentRingCache(10, 10, 1)
	if err != nil {
		t.Fatal(err)
	}

	c.SetTimeFunction(func() int64 {
		return time.Now().UnixNano()
	})

	c.Put("Hi", "Hi")
	d, _ := time.ParseDuration("10ms")
	time.Sleep(d)

	if str, ok := c.Get("Hi"); ok {
		t.Error("key found")
	} else if str != "Hi" {
		t.Error("expired key not returned")
	}
}
