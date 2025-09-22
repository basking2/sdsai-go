package math

import (
	"testing"
)

func TestBicubic(t *testing.T) {
	interp, isValid, err := BicubicInterpolator(
		[]float64{0, 1, 2, 3},
		[]float64{0, 1, 2, 3},
		[][]float64{{0, 1, 2, 3},
			{4, 5, 6, 7},
			{8, 9, 10, 11},
			{12, 13, 14, 15}})

	if err != nil {
		t.Error(err)
	}

	if isValid(0, 0) {
		t.Fatal("Invalid point returned as valid.")
	}

	i, _ := interp(2, 2)

	println(i)
}
