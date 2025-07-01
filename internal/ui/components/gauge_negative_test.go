package components

import (
	"math"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// TestGaugeWithNegativeValues tests the gauge behavior with various negative inputs
// that could potentially cause the strings.Repeat panic
func TestGaugeWithNegativeValues(t *testing.T) {
	testCases := []struct {
		width int
		value float64
		desc  string
	}{
		// Edge cases that could result in negative calculations
		{width: 0, value: 50.0, desc: "zero width"},
		{width: -1, value: 50.0, desc: "negative width"},
		{width: 1, value: -50.0, desc: "negative value"},
		{width: -5, value: -25.0, desc: "both negative"},

		// Extreme floating point values
		{width: 1, value: math.Inf(1), desc: "positive infinity"},
		{width: 1, value: math.Inf(-1), desc: "negative infinity"},
		{width: 1, value: math.NaN(), desc: "NaN value"},

		// Very large values that could overflow
		{width: 1, value: 1e10, desc: "very large value"},
		{width: 1, value: -1e10, desc: "very large negative value"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			// This should not panic regardless of input
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("PANIC with %s: %v", tc.desc, r)
				}
			}()

			gauge := NewGauge(tc.width)
			result := gauge.Render(tc.value, "test")

			// Result should not be empty (we should get something reasonable)
			if result == "" {
				t.Errorf("Empty result for %s", tc.desc)
			}

			// Result should not contain obvious error indicators
			if strings.Contains(result, "panic") || strings.Contains(result, "error") {
				t.Errorf("Result contains error indicators for %s: %s", tc.desc, result)
			}
		})
	}
}

// TestGaugeFilledCalculationEdgeCases tests specific edge cases where
// floating point arithmetic could produce unexpected results
func TestGaugeFilledCalculationEdgeCases(t *testing.T) {
	// Test cases designed to trigger problematic floating point calculations
	testCases := []struct {
		width    int
		value    float64
		expected string // what we expect to NOT see
		desc     string
	}{
		{1, 100.0, "panic", "width=1, 100% - should not panic"},
		{1, 99.99999, "panic", "width=1, near 100% - should not panic"},
		{2, 100.0, "panic", "width=2, 100% - should not panic"},
		{3, 66.66666, "panic", "width=3, repeating decimal - should not panic"},
		{10, 33.33333, "panic", "width=10, repeating decimal - should not panic"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("PANIC: %s caused panic: %v", tc.desc, r)
				}
			}()

			gauge := NewGauge(tc.width)

			// Manually calculate what the old problematic code would do
			filled := int((tc.value / 100.0) * float64(tc.width))
			empty := tc.width - filled

			t.Logf("Testing %s: width=%d, value=%f, calculated filled=%d, empty=%d",
				tc.desc, tc.width, tc.value, filled, empty)

			// This should work without panic
			result := gauge.Render(tc.value, "test")

			if result == "" {
				t.Errorf("Empty result for %s", tc.desc)
			}

			if strings.Contains(strings.ToLower(result), tc.expected) {
				t.Errorf("Result contains unexpected content '%s' for %s: %s",
					tc.expected, tc.desc, result)
			}
		})
	}
}

// TestMultiGaugeEdgeCases tests MultiGauge for similar edge cases
func TestMultiGaugeEdgeCases(t *testing.T) {
	testCases := []struct {
		width int
		total uint64
		desc  string
	}{
		{0, 100, "zero width"},
		{-1, 100, "negative width"},
		{1, 0, "zero total"},
		{5, 0, "zero total with normal width"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("PANIC in MultiGauge %s: %v", tc.desc, r)
				}
			}()

			mg := NewMultiGauge(tc.width)

			// Add some segments
			if tc.total > 0 {
				mg.AddSegment(tc.total/3, "test1", lipgloss.NewStyle())
				mg.AddSegment(tc.total/2, "test2", lipgloss.NewStyle())
			}

			result := mg.Render(tc.total)

			if result == "" && tc.width > 0 {
				t.Errorf("Unexpected empty result for %s", tc.desc)
			}
		})
	}
}
