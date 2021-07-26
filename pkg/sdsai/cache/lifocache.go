package cache

import (
	"container/heap"
	"time"
)

// A heap that expires the first added string first.
//
// This also allows for a string to be refreshed. That is, be put
// at the back of the line.
//
// This cache implementation does not enforce a size limit. It only
// orders items for eviction. The user must evict them to reach a desired
// Len() (size).
type LIFOCache struct {
	// The heap of strings.
	Keys []string

	// The heap of satellite data.
	Items []interface{}

	// The heap of insertion values.
	AddedTime []int64

	// The map of strings to their indexes in the arrays.
	Indexes map[string]int

	// When a string is removed, the eviction handler is called and given the
	// evicted string and associated data.
	//
	// This allows users of this class to have it drive eviction of other resources.
	//
	// The eviction handler is not called when a string is refreshed / re-added.
	EvictionHandlers []func(s string, data interface{})

	// A function that returns the "time" an element is added.
	//
	// This is used to sort the elements for removal where the
	// smallest int64 value is considered first.
	//
	// These integers need not be time, but it is convenient to think of them
	// that way.
	TimeFunction func() int64
}

// Construct a new LIFOCache that uses the system clock in seconds
// to order the strings added.
func NewLIFOCache() *LIFOCache {
	h := LIFOCache{
		Keys:             []string{},
		Items:            []interface{}{},
		Indexes:          make(map[string]int),
		AddedTime:        []int64{},
		EvictionHandlers: []func(string, interface{}){},
		TimeFunction: func() int64 {
			return time.Now().Unix()
		},
	}

	return &h
}

func (c *LIFOCache) Len() int {
	return len(c.Items)
}

func (c *LIFOCache) Less(i, j int) bool {
	return c.AddedTime[i] < c.AddedTime[j]
}

func (c *LIFOCache) Swap(i, j int) {
	// Swap list items.
	c.AddedTime[i], c.AddedTime[j] = c.AddedTime[j], c.AddedTime[i]
	c.Items[i], c.Items[j] = c.Items[j], c.Items[i]
	c.Keys[i], c.Keys[j] = c.Keys[j], c.Keys[i]
	c.EvictionHandlers[i], c.EvictionHandlers[j] = c.EvictionHandlers[j], c.EvictionHandlers[i]

	// Update key to index mapping.
	c.Indexes[c.Keys[i]] = i
	c.Indexes[c.Keys[j]] = j

}

// Push a string key into this cache.
//
// *Do not call this directly.* This is an internal API.
func (c *LIFOCache) Push(x interface{}) {
	switch s := x.(type) {
	default:
		panic("This can only handle strings.")
	case string:
		if _, ok := c.Indexes[s]; ok {
			// Pushing, in this case, is moving the string to the end and
			// updateing the time.
			i := c.Indexes[s]
			c.AddedTime[i] = c.TimeFunction()
			c.Swap(i, len(c.Items)-1)
		} else {
			// Add string.
			c.Indexes[s] = len(c.Keys)
			c.AddedTime = append(c.AddedTime, c.TimeFunction())
			c.Keys = append(c.Keys, s)

			if len(c.AddedTime) != len(c.EvictionHandlers) {
				panic("Use the Put function to add elements to this cache.")
			}

			if len(c.AddedTime) != len(c.Items) {
				panic("Use the Put function to add elements to this cache.")
			}
		}
	}
}

// Return the last user-added object and call the Eviction function.
//
// *Do not call this directly.* This is an internal API.
func (c *LIFOCache) Pop() interface{} {
	// The new length. Also as an index value.
	l := len(c.Items) - 1

	// Record the last item.
	k := c.Keys[l]
	i := c.Items[l]
	e := c.EvictionHandlers[l]

	// Slice the two arrays.
	c.Keys = c.Keys[0:l]
	c.Items = c.Items[0:l]
	c.AddedTime = c.AddedTime[0:l]
	c.EvictionHandlers = c.EvictionHandlers[0:l]

	// Remove the key mapping.
	delete(c.Indexes, k)

	e(k, i)

	// Return the last item.
	return i
}

// Add a key to this cache with a given eviction function.
//
// If the key already exists, the object is updated, the AddedTime is updated
// and a 2-tuple with the previous object and true is returned.
//
// If the key does not already exist, the object is added under that key
// and (nil, false) is returned.
func (c *LIFOCache) PutWithHandler(key string, item interface{}, evictionhandler func(string, interface{})) (interface{}, bool) {

	if _, ok := c.Indexes[key]; ok {
		// Update the time and re-heap.
		i := c.Indexes[key]
		o := c.Items[i]
		c.AddedTime[i] = c.TimeFunction()
		heap.Fix(c, i)

		return o, true
	} else {
		// Push our satellite data first, before the heap data.
		c.EvictionHandlers = append(c.EvictionHandlers, evictionhandler)
		c.Items = append(c.Items, item)
		heap.Push(c, key)

		return nil, false
	}
}

func (c *LIFOCache) Put(key string, item interface{}) (interface{}, bool) {
	return c.PutWithHandler(key, item, func(string, interface{}) {})
}

// Get the user data and the time it was added. The last returned boolean
// indicates if the record was found at all.
//
//     if item, addTime, ok := lifoCache.Get("key"); ok {
//         ...
//     }
func (c *LIFOCache) Get(key string) (interface{}, int64, bool) {
	if i, ok := c.Indexes[key]; ok {
		// If here, the key is i the cache. Now check its validity.
		return c.Items[i], c.AddedTime[i], true
	} else {
		return nil, 0, false
	}
}

// Evict the next item, returning the key and value.
//
// If the cache is empty "" and nil are returned.
func (c *LIFOCache) EvictNext() (string, interface{}) {
	if len(c.Keys) > 0 {
		k := c.Keys[0]
		i := heap.Pop(c)
		return k, i
	} else {
		return "", nil
	}
}

// Evict items that are older than the given tm.
// That is the object's added time is less-than tm.
func (c *LIFOCache) EvictOlderThan(tm int64) {
	for len(c.AddedTime) > 0 && c.AddedTime[0] < tm {
		c.EvictNext()
	}
}

// Set the added time of an item and re-heap it.
func (c *LIFOCache) SetAddedTime(key string, tm int64) {
	if idx, ok := c.Indexes[key]; ok {
		c.AddedTime[idx] = tm
		heap.Fix(c, idx)
	}
}

// Remove the given key from the cache.
func (c *LIFOCache) Remove(key string) (interface{}, bool) {
	if i, ok := c.Indexes[key]; ok {
		obj := c.Items[i]

		lasti := len(c.Items) - 1

		// Put i at the end of the arrays.
		c.Swap(i, lasti)

		// Remove the last element.
		c.Keys = c.Keys[0:lasti]
		c.Items = c.Items[0:lasti]
		c.EvictionHandlers = c.EvictionHandlers[0:lasti]
		c.AddedTime = c.AddedTime[0:lasti]

		delete(c.Indexes, key)

		// Fix i.
		heap.Fix(c, i)

		// Return it.
		return obj, true
	} else {
		return nil, false
	}
}

// Return the next key to be returned by a call to EvictNext().
func (c *LIFOCache) MinKey() string {
	return c.Keys[0]
}

// Return the time of the next key and item to be returned by a call to EvictNext().
func (c *LIFOCache) MinTime() int64 {
	return c.AddedTime[0]
}

// Return next item to be returned by a call to EvictNext().
func (c *LIFOCache) MinItem() interface{} {
	return c.Items[0]
}
