package math

import (
	"errors"
  "sort"
)

// A translation of Apache Commons Math.
// Returns an interpolation function, valid function, and an error if any.
func BicubicInterpolator(
	xval, yval []float64,
	fval [][]float64) (
	func(float64, float64) (float64, error),
	func(float64, float64) bool,
	error) {

	if len(xval) == 0 || len(yval) == 0 || len(fval) == 0 {
		return nil, nil, errors.New("An array is zero length")
	}

	// Check input order.
	for i := 1; i < len(xval); i++ {
		if xval[i] <= xval[i-1] {
			return nil, nil, errors.New("X values are not monotonically increasing.")
		}
	}

	for i := 1; i < len(yval); i++ {
		if yval[i] <= yval[i-1] {
			return nil, nil, errors.New("Y values are not monotonically increasing.")
		}
	}

	dFdX := make([][]float64, len(xval))
	dFdY := make([][]float64, len(xval))
	d2FdXdY := make([][]float64, len(xval))
	for i := 0; i < len(xval); i++ {
		dFdX[i] = make([]float64, len(yval))
		dFdY[i] = make([]float64, len(yval))
		d2FdXdY[i] = make([]float64, len(yval))
	}

	for i := 1; i < len(xval)-1; i++ {
		nI := i + 1
		pI := i - 1

		nX := xval[nI]
		pX := xval[pI]

		deltaX := nX - pX

		for j := 1; j < len(yval)-1; j++ {
			nJ := j + 1
			pJ := j - 1

			nY := yval[nJ]
			pY := yval[pJ]

			deltaY := nY - pY

			dFdX[i][j] = (fval[nI][j] - fval[pI][j]) / deltaX
			dFdY[i][j] = (fval[i][nJ] - fval[i][pJ]) / deltaY

			deltaXY := deltaX * deltaY

			d2FdXdY[j][j] = (fval[nI][nJ] - fval[nI][pJ] - fval[pI][nJ] + fval[pI][pJ]) / deltaXY
		}
	}

	isValidFn := func(x, y float64) bool {
		if x < xval[1] || x > xval[len(xval)-2] || y < yval[1] || y > yval[len(yval)-2] {
			// println("----------------")
			// println(x, xval[1], xval[len(xval)-2], y, yval[1], yval[len(yval)-2])
			// println("----------------")
			return false
		} else {
			return true
		}
	}

	interpolateFn, err := bicubicInterplationFunction(xval, yval, fval, dFdX, dFdY, d2FdXdY)
	if err != nil {
		return nil, nil, err
	}

	return interpolateFn, isValidFn, nil
}

func bicubicInterplationFunction(
	xval, yval []float64,
	f, dFdX, dFdY, d2FdXdY [][]float64) (func(x, y float64) (float64, error), error) {

	if len(xval) == 0 || len(yval) == 0 || len(f) == 0 || len(f[0]) == 0 {
		return nil, errors.New("Input length is zero.")
	}
	if len(xval) != len(f) {
		return nil, errors.New("X len does not equal f len.")
	}
	if len(xval) != len(dFdX) {
		return nil, errors.New("Dimensions don't match between x and dFDx.")
	}
	if len(xval) != len(dFdY) {
		return nil, errors.New("Dimensions don't match between x and dFDY.")
	}
	if len(xval) != len(d2FdXdY) {
		return nil, errors.New("Dimensions don't match between x and d2FdXdY.")
	}

	// We skip x and y order checks because it is done in the calling function.

	lastI := len(xval) - 1
	lastJ := len(yval) - 1
	splines := make([][]func(float64, float64) (float64, error), lastI)

	for i := 0; i < lastI; i++ {
		splines[i] = make([]func(float64, float64) (float64, error), lastJ)
		if len(f[i]) != len(yval) {
			return nil, errors.New("Dimension mismatch of f.")
		}
		if len(dFdX[i]) != len(yval) {
			return nil, errors.New("Dimension mismatch of dFdX.")
		}
		if len(dFdY[i]) != len(yval) {
			return nil, errors.New("Dimension mismatch of dFdY.")
		}
		if len(d2FdXdY[i]) != len(yval) {
			return nil, errors.New("Dimension mismatch of d2FdXdY[i].")
		}
		ip1 := i + 1
		xR := xval[ip1] - xval[i]
		for j := 0; j < lastJ; j++ {
			jp1 := j + 1
			yR := yval[jp1] - yval[j]
			xRyR := xR * yR
			beta := []float64{
				f[i][j], f[ip1][j], f[i][jp1], f[ip1][jp1],
				dFdX[i][j] * xR, dFdX[ip1][j] * xR, dFdX[i][jp1] * xR, dFdX[ip1][jp1] * xR,
				dFdY[i][j] * yR, dFdY[ip1][j] * yR, dFdY[i][jp1] * yR, dFdY[ip1][jp1] * yR,
				d2FdXdY[i][j] * xRyR, d2FdXdY[ip1][j] * xRyR, d2FdXdY[i][jp1] * xRyR, d2FdXdY[ip1][jp1] * xRyR}

			splines[i][j] = bicubicFunction(computeSplineCoefficients(beta))
		}
	}

  interpF := func(x, y float64)(float64, error) {
    i := searchIndex(x, xval)
    j := searchIndex(y, yval)

    xN := (x - xval[i]) / (xval[i+1] - xval[i])
    yN := (y - yval[j]) / (yval[j+1] - yval[j])

    return splines[i][j](xN, yN)
  }

	return interpF, nil
}

func computeSplineCoefficients(beta []float64) []float64 {
	NUM_COEFF := 16
	AINV := [][]float64{
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{-3, 3, 0, 0, -2, -1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{2, -2, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, -3, 3, 0, 0, -2, -1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 2, -2, 0, 0, 1, 1, 0, 0},
		{-3, 0, 3, 0, 0, 0, 0, 0, -2, 0, -1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, -3, 0, 3, 0, 0, 0, 0, 0, -2, 0, -1, 0},
		{9, -9, -9, 9, 6, 3, -6, -3, 6, -6, 3, -3, 4, 2, 2, 1},
		{-6, 6, 6, -6, -3, -3, 3, 3, -4, 4, -2, 2, -2, -2, -1, -1},
		{2, 0, -2, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 2, 0, -2, 0, 0, 0, 0, 0, 1, 0, 1, 0},
		{-6, 6, 6, -6, -4, -2, 4, 2, -3, 3, -3, 3, -2, -1, -2, -1},
		{4, -4, -4, 4, 2, 2, -2, -2, 2, -2, 2, -2, 1, 1, 1, 1}}

	a := make([]float64, NUM_COEFF)
	for i := 0; i < NUM_COEFF; i++ {
		result := float64(0)
		row := AINV[i]
		for j := 0; j < NUM_COEFF; j++ {
			result +=  (row[j] * beta[j])
		}
		a[i] = result
	}
	return a
}

func bicubicFunction(coeff []float64) func(float64, float64) (float64, error) {

	N := 4

	a := make([][]float64, N)
	for j := range a {
		a[j] = make([]float64, N)
		aJ := a[j]
		for i := range aJ {
			aJ[i] = coeff[i*N+j]
		}
	}

	applyFn := func(pX, pY []float64, coeff [][]float64) float64 {
		result := float64(0)
		for i := 0; i < N; i++ {
			r := LinearCombination(coeff[i], pY)
			result += r * pX[i]
		}
		return result
	}

	return func(x, y float64) (float64, error) {
		if x < 0 || x > 1 {
			return 0, errors.New("Out of range x 0 - 1.")
		}
		if y < 0 || y > 1 {
			return 0, errors.New("Out of range y 0 - 1.")
		}

		x2 := x * x
		x3 := x2 * x
		pX := []float64{1, x, x2, x3}

		y2 := y * y
		y3 := y2 * y
		pY := []float64{1, y, y2, y3}

		return applyFn(pX, pY, a), nil
	}
}

// Search a sorted set of values for c.
func searchIndex(c float64, val []float64) int {

  idx := sort.SearchFloat64s(val, c)

  // Too low or high. Return -1.
  if idx == 0 || idx >= len(val) {
    return -1
  }

  // The returned index is used as a rage. Ensure that idx+1 is in the array list.
  return idx - 1
}

func LinearCombination(a, b []float64) float64 {
	result := float64(0)
	for i := range a {
		result = result + a[i] * b[i]
	}
	return result
}
