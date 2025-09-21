package consistenthash

import (
	"hash/fnv"
)

// HashToInt takes an input string and a modulus and returns an integer in [0, mod).
// If mod <= 0 it returns 0.
func HashToInt(s string, mod int) int {
	if mod <= 0 {
		return 0
	}
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return int(h.Sum64() % uint64(mod))
}
