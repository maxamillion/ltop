package utils

import (
	"math"
)

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MaxFloat(a, b float64) float64 {
	return math.Max(a, b)
}

func MinFloat(a, b float64) float64 {
	return math.Min(a, b)
}

func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func ClampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func AbsFloat(x float64) float64 {
	return math.Abs(x)
}

func Round(x float64) int {
	return int(math.Round(x))
}

func RoundToDecimal(x float64, decimals int) float64 {
	multiplier := math.Pow(10, float64(decimals))
	return math.Round(x*multiplier) / multiplier
}

func Percentage(part, total uint64) float64 {
	if total == 0 {
		return 0.0
	}
	return float64(part) / float64(total) * 100.0
}

func Average(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}

	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func Sum(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum
}

func SumInt(values []int) int {
	sum := 0
	for _, v := range values {
		sum += v
	}
	return sum
}

func SumUint64(values []uint64) uint64 {
	var sum uint64
	for _, v := range values {
		sum += v
	}
	return sum
}

func SafeDivide(numerator, denominator float64) float64 {
	if denominator == 0 {
		return 0.0
	}
	return numerator / denominator
}

func SafeDivideUint64(numerator, denominator uint64) float64 {
	if denominator == 0 {
		return 0.0
	}
	return float64(numerator) / float64(denominator)
}

func LinearInterpolate(x, x1, y1, x2, y2 float64) float64 {
	if x2 == x1 {
		return y1
	}
	return y1 + (x-x1)*(y2-y1)/(x2-x1)
}

func MapRange(value, fromMin, fromMax, toMin, toMax float64) float64 {
	if fromMax == fromMin {
		return toMin
	}
	return toMin + (value-fromMin)*(toMax-toMin)/(fromMax-fromMin)
}
