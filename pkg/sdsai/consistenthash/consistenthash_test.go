package consistenthash

import (
	"fmt"
	"testing"
)

func TestHashDeterministic(t *testing.T) {
	mod := 1000
	keys := []string{
		"apple", "banana", "orange", "key-1", "key-2", "",
	}
	for _, k := range keys {
		a := HashToInt(k, mod)
		b := HashToInt(k, mod)
		if a != b {
			t.Fatalf("HashToInt not deterministic for %q: %d != %d", k, a, b)
		}
	}
}

func TestHashRangeAndModulo(t *testing.T) {
	tests := []int{1, 2, 10, 100, 1024}
	for _, mod := range tests {
		for i := 0; i < 100; i++ {
			k := fmt.Sprintf("k-%d", i)
			v := HashToInt(k, mod)
			if v < 0 || v >= mod {
				t.Fatalf("result out of range for mod=%d, key=%q: %d", mod, k, v)
			}
		}
	}
}

func TestModZeroOrNegative(t *testing.T) {
	if got := HashToInt("anything", 0); got != 0 {
		t.Fatalf("expected 0 for mod=0, got %d", got)
	}
	if got := HashToInt("anything", -5); got != 0 {
		t.Fatalf("expected 0 for negative mod, got %d", got)
	}
}

func TestDistributionBasic(t *testing.T) {
	mod := 10
	counts := make([]int, mod)
	N := 1000
	for i := 0; i < N; i++ {
		k := fmt.Sprintf("item-%d", i)
		counts[HashToInt(k, mod)]++
	}
	// Ensure each bucket got at least one item (very likely with N=1000 and mod=10).
	for i, c := range counts {
		if c == 0 {
			t.Fatalf("bucket %d received zero items (counts: %v)", i, counts)
		}
	}
}
