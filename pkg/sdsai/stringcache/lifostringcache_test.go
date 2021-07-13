package stringcache

import (
	"container/heap"
	"testing"
	"time"
)

func TestLIFOStringCache(t *testing.T) {

	cache := NewLIFOStringCache()

	t.Run("Add", func(t *testing.T) {
		heap.Push(cache, "string 1")
		heap.Push(cache, "string 2")

		d, _ := time.ParseDuration("1s")
		time.Sleep(d)
		heap.Push(cache, "string 3")
		time.Sleep(d)

		if cache.Len() != 3 {
			t.Errorf("Heap length should be 3 but was %d.", cache.Len())
		}
	})

	t.Run("Remove and Refresh", func(t *testing.T) {
		switch v := heap.Pop(cache).(type) {
		default:
			t.Error("Heap returned non-string type.")
		case string:
			if v == "string 1" {
				heap.Push(cache, "string 2")
			} else if v == "string 2" {
				heap.Push(cache, "string 1")
			} else {
				t.Errorf("String 1 or 2 was not returned.")
			}
		}
	})

	t.Run("Remove string 1 or 2", func(t *testing.T) {
		switch v := heap.Pop(cache).(type) {
		default:
			t.Error("Heap returned non-string type.")
		case string:
			if v != "string 3" {
				t.Errorf("String 1 or 2 was not returned.")
			}
		}
	})

}
