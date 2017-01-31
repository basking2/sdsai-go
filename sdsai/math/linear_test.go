
package math


import (
	"testing"
)

func TestLinearInterpolator(t *testing.T) {
	li,err := LinearInterpolator([]float64{0, 1}, []float64{3, 4})
	if err != nil {
		t.Fatal(err)
	}

	v, err := li(0.5)
	if err != nil {
		t.Fatal(err)
	}

	println(v)
}
