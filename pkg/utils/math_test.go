package utils

import "testing"

func TestMax(t *testing.T) {
	testCases := []struct {
		a, b     int
		expected int
	}{
		{5, 10, 10},
		{10, 5, 10},
		{0, 0, 0},
		{-5, 5, 5},
		{-10, -5, -5},
	}

	for _, tc := range testCases {
		result := Max(tc.a, tc.b)
		if result != tc.expected {
			t.Errorf("Max(%d, %d) = %d; expected %d", tc.a, tc.b, result, tc.expected)
		}
	}
}

func TestMin(t *testing.T) {
	testCases := []struct {
		a, b     int
		expected int
	}{
		{5, 10, 5},
		{10, 5, 5},
		{0, 0, 0},
		{-5, 5, -5},
		{-10, -5, -10},
	}

	for _, tc := range testCases {
		result := Min(tc.a, tc.b)
		if result != tc.expected {
			t.Errorf("Min(%d, %d) = %d; expected %d", tc.a, tc.b, result, tc.expected)
		}
	}
}

func TestClamp(t *testing.T) {
	testCases := []struct {
		value, min, max int
		expected        int
	}{
		{5, 0, 10, 5},
		{-5, 0, 10, 0},
		{15, 0, 10, 10},
		{5, 5, 10, 5},
		{10, 5, 10, 10},
	}

	for _, tc := range testCases {
		result := Clamp(tc.value, tc.min, tc.max)
		if result != tc.expected {
			t.Errorf("Clamp(%d, %d, %d) = %d; expected %d", tc.value, tc.min, tc.max, result, tc.expected)
		}
	}
}

func TestPercentage(t *testing.T) {
	testCases := []struct {
		part, total uint64
		expected    float64
	}{
		{50, 100, 50.0},
		{25, 100, 25.0},
		{0, 100, 0.0},
		{100, 100, 100.0},
		{10, 0, 0.0}, // Division by zero case
	}

	for _, tc := range testCases {
		result := Percentage(tc.part, tc.total)
		if result != tc.expected {
			t.Errorf("Percentage(%d, %d) = %f; expected %f", tc.part, tc.total, result, tc.expected)
		}
	}
}

func TestAverage(t *testing.T) {
	testCases := []struct {
		values   []float64
		expected float64
	}{
		{[]float64{1, 2, 3, 4, 5}, 3.0},
		{[]float64{10, 20}, 15.0},
		{[]float64{5}, 5.0},
		{[]float64{}, 0.0}, // Empty slice case
	}

	for _, tc := range testCases {
		result := Average(tc.values)
		if result != tc.expected {
			t.Errorf("Average(%v) = %f; expected %f", tc.values, result, tc.expected)
		}
	}
}

func TestSafeDivide(t *testing.T) {
	testCases := []struct {
		numerator, denominator float64
		expected               float64
	}{
		{10, 2, 5.0},
		{0, 5, 0.0},
		{10, 0, 0.0}, // Division by zero case
		{-10, 2, -5.0},
	}

	for _, tc := range testCases {
		result := SafeDivide(tc.numerator, tc.denominator)
		if result != tc.expected {
			t.Errorf("SafeDivide(%f, %f) = %f; expected %f", tc.numerator, tc.denominator, result, tc.expected)
		}
	}
}
