package math

import (
	"math"
	"testing"
)

func TestLinearInterpolator(t *testing.T) {
	li, err := LinearInterpolator([]float64{0, 1}, []float64{0, 1})
	if err != nil {
		t.Fatal(err)
	}

	v, err := li(0.5)
	if err != nil {
		t.Fatal(err)
	}

	if math.Abs(v-0.5) > 000.1 {
		t.Fatal("Expected near 0.5 but got ", v)
	}
}
