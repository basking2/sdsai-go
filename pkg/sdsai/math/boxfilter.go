package math

import (
	gomath "math"
)

type BoxFilter struct {
	size   int
	before int
	after  int
}

func NewBoxFilter() BoxFilter {
	return BoxFilter{3, 1, 1}
}

func NewCustomBoxFilter(before, after int) BoxFilter {
	return BoxFilter{before + after + 1, before, after}
}

// Implement a 2-pass box blur.
// 2-pass does more work but possibly makes better use of caching.
func (bf BoxFilter) Filter64(data []float64, height, width int) []float64 {
	result := [2][]float64{
		make([]float64, len(data)),
		make([]float64, len(data))}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var start int
			var stop int
			if x-bf.before < 0 {
				start = 0
			} else {
				start = x - bf.before
			}

			if x+bf.after < width {
				stop = x + bf.after + 1
			} else {
				stop = width
			}

			count := 0
			value := float64(0)

			for i := start; i < stop; i++ {
				tmpValue := float64(data[i+y*width])
				if !gomath.IsNaN(tmpValue) {
					value += tmpValue
					count += 1
				}
			}

			if count > (bf.before+bf.after+1)/2 {
				result[0][x+y*width] = value / float64(count)
			} else {
				result[0][x+y*width] = data[x+y*width]
			}
		}
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var start int
			var stop int
			if y-bf.before < 0 {
				start = 0
			} else {
				start = y - bf.before
			}

			if y+bf.after < height {
				stop = y + bf.after + 1
			} else {
				stop = height
			}

			count := 0
			value := float64(0)

			for i := start; i < stop; i++ {
				tmpValue := result[0][x+i*width]
				if !gomath.IsNaN(tmpValue) {
					value += tmpValue
					count += 1
				}
			}

			if count > (bf.before+bf.after+1)/2 {
				result[1][x+y*width] = value / float64(count)
			} else {
				result[1][x+y*width] = gomath.NaN()
			}
		}
	}

	return result[1]
}

// Implement a 2-pass box blur.
// 2-pass does more work but possibly makes better use of caching.
func (bf BoxFilter) Filter32(data []float32, height, width int) []float32 {
	result := [2][]float32{
		make([]float32, len(data)),
		make([]float32, len(data))}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var start int
			var stop int
			if x-bf.before < 0 {
				start = 0
			} else {
				start = x - bf.before
			}

			if x+bf.after < width {
				stop = x + bf.after + 1
			} else {
				stop = width
			}

			count := 0
			value := float32(0)

			for i := start; i < stop; i++ {
				tmpValue := float32(data[i+y*width])
				if !gomath.IsNaN(float64(tmpValue)) {
					value += tmpValue
					count += 1
				}
			}

			if count > (bf.before+bf.after+1)/2 {
				result[0][x+y*width] = value / float32(count)
			} else {
				result[0][x+y*width] = data[x+y*width]
			}
		}
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var start int
			var stop int
			if y-bf.before < 0 {
				start = 0
			} else {
				start = y - bf.before
			}

			if y+bf.after < height {
				stop = y + bf.after + 1
			} else {
				stop = height
			}

			count := 0
			value := float32(0)

			for i := start; i < stop; i++ {
				tmpValue := result[0][x+i*width]
				if !gomath.IsNaN(float64(tmpValue)) {
					value += tmpValue
					count += 1
				}
			}

			if count > (bf.before+bf.after+1)/2 {
				result[1][x+y*width] = value / float32(count)
			} else {
				result[1][x+y*width] = float32(gomath.NaN())
			}
		}
	}

	return result[1]
}
