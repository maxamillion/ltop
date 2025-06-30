package components

import (
	"fmt"
	"strings"
	"testing"
	
	"github.com/charmbracelet/lipgloss"
)

// TestGaugeGo124Regression specifically tests the issue reported on Go 1.24.4
// where strings.Repeat was being called with negative values causing panics
func TestGaugeGo124Regression(t *testing.T) {
	// Test the exact conditions that caused the panic on Go 1.24.4
	// Storage view creates gauges with width=30 and renders filesystem usage
	
	// Create a gauge exactly like the storage view does
	gauge := NewGauge(30)
	
	// Test various percentage values that could cause issues
	problemValues := []float64{
		// Values that could cause floating point precision issues
		99.9999999,  // Very close to 100%
		100.0,       // Exactly 100%
		66.66666667, // Repeating decimal
		33.33333333, // Repeating decimal  
		95.1,        // Decimal value
		75.25,       // Quarter precision
		
		// Edge cases
		0.1,         // Very small value
		99.1,        // Near 100%
		100.1,       // Slightly over 100% (should be clamped)
	}
	
	for _, value := range problemValues {
		t.Run(fmt.Sprintf("value_%.9f", value), func(t *testing.T) {
			// This is the exact call pattern from storage view
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("PANIC on Go 1.24.4 scenario with value %.9f: %v", value, r)
				}
			}()
			
			// Test the exact call pattern from storage.go:53
			// gauge := sv.fsGauges[i].Render(filesystem.UsedPercent, mountpoint)
			result := gauge.Render(value, "/")
			
			if result == "" {
				t.Errorf("Empty result for value %.9f", value)
			}
			
			// Should not contain panic-related strings
			if strings.Contains(result, "panic") {
				t.Errorf("Result contains panic string for value %.9f: %s", value, result)
			}
			
			// Log what the problematic calculation would have been
			filled := int((value / 100.0) * float64(30))
			empty := 30 - filled
			
			if filled < 0 || empty < 0 {
				t.Logf("Go 1.24.4 REGRESSION: value=%.9f would cause panic - filled=%d, empty=%d", 
					value, filled, empty)
			}
		})
	}
}

// TestGaugeStringRepeatSafety ensures all calls to strings.Repeat have non-negative arguments
func TestGaugeStringRepeatSafety(t *testing.T) {
	// Test all the ways strings.Repeat is called in gauge.go
	
	// Test regular Render method
	gauge := NewGauge(10)
	for value := -50.0; value <= 150.0; value += 0.1 {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("PANIC in Render with value %.1f: %v", value, r)
				}
			}()
			gauge.Render(value, "test")
		}()
	}
	
	// Test RenderWithColors method
	lowStyle := lipgloss.NewStyle()
	midStyle := lipgloss.NewStyle()
	highStyle := lipgloss.NewStyle()
	
	for value := -50.0; value <= 150.0; value += 0.1 {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("PANIC in RenderWithColors with value %.1f: %v", value, r)
				}
			}()
			gauge.RenderWithColors(value, "test", lowStyle, midStyle, highStyle)
		}()
	}
	
	// Test MultiGauge with various scenarios
	mg := NewMultiGauge(10)
	mg.AddSegment(50, "test1", lipgloss.NewStyle())
	mg.AddSegment(30, "test2", lipgloss.NewStyle())
	
	// Test with various total values including edge cases
	totals := []uint64{0, 1, 100, 1000, ^uint64(0)} // including max uint64
	
	for _, total := range totals {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("PANIC in MultiGauge.Render with total %d: %v", total, r)
				}
			}()
			mg.Render(total)
		}()
	}
}

// TestGaugeExtremeWidths tests gauge behavior with extreme width values
func TestGaugeExtremeWidths(t *testing.T) {
	extremeWidths := []int{
		0, -1, -100,    // Negative and zero widths
		1, 2,           // Minimal widths  
		1000, 10000,    // Large widths
	}
	
	for _, width := range extremeWidths {
		t.Run(fmt.Sprintf("width_%d", width), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("PANIC with extreme width %d: %v", width, r)
				}
			}()
			
			gauge := NewGauge(width)
			
			// Test with a value that typically causes issues
			result := gauge.Render(100.0, "test")
			
			if result == "" && gauge.Width > 0 {
				t.Errorf("Unexpected empty result for width %d", width)
			}
		})
	}
}

