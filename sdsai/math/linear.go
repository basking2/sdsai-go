package math

import (
	"errors"
)

// Polynomial function evaluation.
func evaluation(coefficients []float64, arg float64) float64 {
	result := coefficients[len(coefficients)-1]
	for j := len(coefficients) - 2; j >=0; j-- {
		result = arg * result + coefficients[j]
	}

	return result
}

func LinearInterpolator(x, y []float64) (func(float64)(float64, error), error) {
	if len(x) != len(y) {
		return nil, errors.New("x and y dimensions must be the same length.")
	}

	intervals := len(x) - 1

	m := make([]float64, intervals)

	for i := 0; i < intervals; i++ {
		m[i] = (y[i+1] - y[i]) / (x[i+1] - x[i])
	}

	polynomials := make([]func(float64)float64, intervals)
	for i := 0; i < intervals; i++ {
		polynomials[i] = func(arg float64) float64 {
			return evaluation([]float64{ y[i], m[i] }, arg)
		}
	}

	return func(arg float64) (float64, error) {
		if arg < x[0] || arg > x[len(x)-1] {
			return 0, errors.New("Out of range.")
		}

		// Find interval.
		i := searchIndex(arg, x)

		return polynomials[i](arg - x[i]), nil
	}, nil
}

