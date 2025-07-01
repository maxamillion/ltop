package components

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestNewGauge(t *testing.T) {
	// Test normal width
	gauge := NewGauge(50)
	if gauge == nil {
		t.Fatal("NewGauge returned nil")
	}
	if gauge.Width != 50 {
		t.Errorf("Expected width 50, got %d", gauge.Width)
	}

	// Test zero width protection
	gaugeZero := NewGauge(0)
	if gaugeZero.Width != 1 {
		t.Errorf("Expected width 1 for zero input, got %d", gaugeZero.Width)
	}

	// Test negative width protection
	gaugeNeg := NewGauge(-5)
	if gaugeNeg.Width != 1 {
		t.Errorf("Expected width 1 for negative input, got %d", gaugeNeg.Width)
	}
}

func TestGaugeRender(t *testing.T) {
	gauge := NewGauge(10)

	// Test normal values
	result := gauge.Render(50.0, "Test")
	if result == "" {
		t.Error("Gauge render should not return empty string")
	}

	// Test edge cases that could cause panic
	testCases := []struct {
		value float64
		width int
		name  string
	}{
		{100.0, 1, "100% with width 1"},
		{99.9, 1, "99.9% with width 1"},
		{100.0, 2, "100% with width 2"},
		{50.0, 3, "50% with width 3"},
		{33.33, 3, "33.33% with width 3 (floating point precision)"},
		{66.67, 3, "66.67% with width 3 (floating point precision)"},
		{0.0, 1, "0% with width 1"},
		{-10.0, 5, "negative value"},
		{110.0, 5, "over 100%"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGauge(tc.width)
			// This should not panic
			result := g.Render(tc.value, "test")
			if result == "" {
				t.Errorf("Render returned empty string for case: %s", tc.name)
			}
		})
	}
}

func TestGaugeRenderNoPanic(t *testing.T) {
	// This test specifically targets the original panic condition
	// where filled > width could cause empty to be negative

	// Create scenarios that could trigger the bug
	testCases := []struct {
		width int
		value float64
		name  string
	}{
		{1, 100.0, "width=1, value=100%"},
		{2, 100.0, "width=2, value=100%"},
		{3, 99.9, "width=3, value=99.9%"},
		{4, 75.1, "width=4, value=75.1%"},
		{5, 60.1, "width=5, value=60.1%"},
		{10, 33.4, "width=10, value=33.4%"},
		{20, 95.1, "width=20, value=95.1%"},
		{30, 100.0, "width=30, value=100%"}, // This matches the storage view usage
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gauge := NewGauge(tc.width)

			// Calculate what the original buggy code would do
			filled := int((tc.value / 100.0) * float64(tc.width))
			empty := tc.width - filled

			// This should not panic even if empty would be negative
			result := gauge.Render(tc.value, "test")

			// Verify the result makes sense
			if result == "" {
				t.Errorf("Render returned empty string for %s", tc.name)
			}

			// Check that the result doesn't contain negative repeat artifacts
			if strings.Contains(result, "strings.Repeat") {
				t.Errorf("Result contains error text, likely from panic recovery: %s", result)
			}

			// If this was the problematic case, log it for debugging
			if empty < 0 {
				t.Logf("Case %s would have caused panic: filled=%d, empty=%d, width=%d",
					tc.name, filled, empty, tc.width)
			}
		})
	}
}

func TestGaugeRenderWithColors(t *testing.T) {
	gauge := NewGauge(10)

	lowColor := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	midColor := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00"))
	highColor := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))

	// Test the same edge cases for RenderWithColors
	testCases := []float64{0.0, 33.33, 50.0, 66.67, 99.9, 100.0, -10.0, 110.0}

	for _, value := range testCases {
		t.Run(fmt.Sprintf("value_%.1f", value), func(t *testing.T) {
			// This should not panic
			result := gauge.RenderWithColors(value, "test", lowColor, midColor, highColor)
			if result == "" {
				t.Errorf("RenderWithColors returned empty string for value %.1f", value)
			}
		})
	}
}

func TestNewMultiGauge(t *testing.T) {
	// Test normal width
	mg := NewMultiGauge(50)
	if mg == nil {
		t.Fatal("NewMultiGauge returned nil")
	}
	if mg.Width != 50 {
		t.Errorf("Expected width 50, got %d", mg.Width)
	}

	// Test zero width protection
	mgZero := NewMultiGauge(0)
	if mgZero.Width != 1 {
		t.Errorf("Expected width 1 for zero input, got %d", mgZero.Width)
	}

	// Test negative width protection
	mgNeg := NewMultiGauge(-5)
	if mgNeg.Width != 1 {
		t.Errorf("Expected width 1 for negative input, got %d", mgNeg.Width)
	}
}

func TestMultiGaugeRender(t *testing.T) {
	mg := NewMultiGauge(10)

	// Test empty gauge (total=0 should return empty bar characters)
	result := mg.Render(0)
	// Check that it contains the expected number of empty characters
	expectedChars := strings.Repeat("░", 10)
	if result != expectedChars {
		t.Errorf("Expected empty gauge to be %q, got %q", expectedChars, result)
	}

	// Test with segments
	color := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	mg.AddSegment(30, "test", color)
	mg.AddSegment(20, "test2", color)

	result = mg.Render(100)
	if result == "" {
		t.Error("MultiGauge render should not return empty string")
	}
}
