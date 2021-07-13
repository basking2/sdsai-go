package stringcache

import (
	"testing"
)

func TestLIFOCache(t *testing.T) {

	cache := NewLIFOCache()
	cacheItemClock := int64(0)
	cache.TimeFunction = func() int64 {
		cacheItemClock++
		return cacheItemClock
	}

	string1removed := false
	string2removed := false
	string3removed := false

	t.Run("Add", func(t *testing.T) {

		cache.PutWithHandler("string 1", "string 1 obj", func(string, interface{}) {
			string1removed = true
		})
		cache.PutWithHandler("string 2", "string 2 obj", func(string, interface{}) {
			string2removed = true
		})

		cache.PutWithHandler("string 3", "string 3 obj", func(string, interface{}) {
			string3removed = true
		})

		if cache.Len() != 3 {
			t.Errorf("Heap length should be 3 but was %d.", cache.Len())
		}
	})

	t.Run("Remove and Refresh", func(t *testing.T) {

		k, v := cache.EvictNext()

		if k != "string 1" && k != "string 2" {
			t.Errorf("Neither string 1 nor 2 but %s was returned.", k)
		}

		switch vstr := v.(type) {
		case string:
			if vstr != "string 1 obj" && vstr != "string 2 obj" {
				t.Errorf("Neither string 1 obj nor 2 obj but %s was returned.", vstr)
			}
		default:
			t.Errorf("Non-string object returned: %s", v)
		}

		if string1removed {
			cache.Put("string 2", "string 2 obj")
		} else if string2removed {
			cache.Put("string 1", "string 1 obj")
		} else {
			t.Errorf("String 1 or 2 was not returned.")
		}
	})

	t.Run("Remove string 1 or 2", func(t *testing.T) {
		k, v := cache.EvictNext()

		if !string3removed {
			t.Errorf("String 3 was not removed.")
		}

		if k != "string 3" {
			t.Errorf("String 3 was not returned but %s.", k)
		}

		switch vstr := v.(type) {
		default:
			t.Errorf("Non-string object returned: %s", v)
		case string:
			if vstr != "string 3 obj" {
				t.Errorf("Not string 3 obj but %s was returned.", vstr)
			}
		}
	})

}
