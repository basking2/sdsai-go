package stringcache

import "time"

// A heap that expires the first added string first.
//
// This also allows for a string to be refreshed. That is, be put
// at the back of the line.
type LIFOStringCache struct {
	// The heap of strings.
	Items []string

	// The heap of insertion values.
	AddedTime []int64

	// The map of strings to their indexes in the arrays.
	Indexes map[string]int

	// A function that returns the "time" an element is added.
	//
	// This is used to sort the elements for removal where the
	// smallest int64 value is considered first.
	//
	// These integers need not be time, but it is convenient to think of them
	// that way.
	TimeFunction func() int64
}

// Construct a new LIFOStringCache that uses the system clock in seconds
// to order the strings added.
func NewLIFOStringCache() *LIFOStringCache {
	h := LIFOStringCache{
		Items:     []string{},
		Indexes:   make(map[string]int),
		AddedTime: []int64{},
		TimeFunction: func() int64 {
			return time.Now().Unix()
		},
	}

	return &h
}

func (c *LIFOStringCache) Len() int {
	return len(c.Items)
}

func (c *LIFOStringCache) Less(i, j int) bool {
	return c.AddedTime[i] < c.AddedTime[j]
}

func (c *LIFOStringCache) Swap(i, j int) {
	// Swap list items.
	c.AddedTime[i], c.AddedTime[j] = c.AddedTime[j], c.AddedTime[i]
	c.Items[i], c.Items[j] = c.Items[j], c.Items[i]

	// Update key to index mapping.
	c.Indexes[c.Items[i]] = i
	c.Indexes[c.Items[j]] = j
}

func (c *LIFOStringCache) Push(x interface{}) {
	switch s := x.(type) {
	default:
		panic("This can only handle strings.")
	case string:
		if _, ok := c.Indexes[s]; ok {
			// Pushing, in this case, is moving the string to the end.
			i := c.Indexes[s]
			c.AddedTime[i] = c.TimeFunction()
			c.Swap(i, len(c.Items)-1)
		} else {
			// Add string.
			c.Indexes[s] = len(c.Items)
			c.AddedTime = append(c.AddedTime, c.TimeFunction())
			c.Items = append(c.Items, s)
		}
	}
}

func (c *LIFOStringCache) Pop() interface{} {
	// The new length. Also as an index value.
	l := len(c.Items) - 1

	// Record the last item.
	i := c.Items[l]

	// Slice the two arrays.
	c.Items = c.Items[0:l]
	c.AddedTime = c.AddedTime[0:l]

	// Remove the key mapping.
	delete(c.Indexes, i)

	// Return the last item.
	return i
}
