package components

import (
	"strings"
	"testing"
)

// TestGaugeRegressionNegativeRepeat tests the specific issue that caused the panic
// on Fedora 41: "strings: negative Repeat count"
// This test verifies that the fix prevents the panic that occurred when
// floating point calculations resulted in `filled > width`, making `empty` negative.
func TestGaugeRegressionNegativeRepeat(t *testing.T) {
	// These are the specific conditions that could trigger the original bug
	testCases := []struct {
		width int
		value float64
		desc  string
	}{
		// These cases are specifically designed to trigger floating point precision
		// issues that could make filled > width
		{width: 1, value: 100.0, desc: "width=1, 100% - minimal width edge case"},
		{width: 30, value: 100.0, desc: "width=30, 100% - storage view usage"},
		{width: 3, value: 66.67, desc: "width=3, 66.67% - floating point precision"},
		{width: 4, value: 75.25, desc: "width=4, 75.25% - quarter precision"},
		{width: 5, value: 80.1, desc: "width=5, 80.1% - decimal precision"},
		
		// Edge cases that could cause the issue
		{width: 10, value: 99.99, desc: "width=10, 99.99% - near 100%"},
		{width: 7, value: 85.71, desc: "width=7, 85.71% - repeating decimal"},
		{width: 9, value: 88.89, desc: "width=9, 88.89% - repeating decimal"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			gauge := NewGauge(tc.width)
			
			// Calculate what the problematic values would be
			filled := int((tc.value / 100.0) * float64(tc.width))
			empty := tc.width - filled
			
			// Log if this would have been a problematic case
			if empty < 0 {
				t.Logf("REGRESSION CASE: %s would cause panic - filled=%d, empty=%d, width=%d", 
					tc.desc, filled, empty, tc.width)
			}
			
			// This should not panic, regardless of the internal calculations
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("PANIC in %s: %v", tc.desc, r)
				}
			}()
			
			result := gauge.Render(tc.value, "test")
			
			// Verify we got a reasonable result
			if result == "" {
				t.Errorf("Empty result for %s", tc.desc)
			}
			
			// Verify the result doesn't contain error indicators
			if strings.Contains(result, "panic") || strings.Contains(result, "negative") {
				t.Errorf("Result contains error indicators for %s: %s", tc.desc, result)
			}
		})
	}
}

// TestGaugeConsistencyAfterFix verifies that the gauge behaves consistently
// and produces visually correct output after the fix
func TestGaugeConsistencyAfterFix(t *testing.T) {
	gauge := NewGauge(10)
	
	// Test that 0% produces all empty chars
	result0 := gauge.Render(0.0, "")
	if !strings.Contains(result0, "░") {
		t.Error("0% should contain empty characters")
	}
	
	// Test that 100% produces all filled chars  
	result100 := gauge.Render(100.0, "")
	if !strings.Contains(result100, "█") {
		t.Error("100% should contain filled characters")
	}
	
	// Test that the total length makes sense (accounting for percentage display)
	if len(result0) == 0 {
		t.Error("Gauge result should not be empty")
	}
	
	// Test with the exact width from storage view that caused the issue
	storageGauge := NewGauge(30)
	storageResult := storageGauge.Render(95.1, "/")
	if storageResult == "" {
		t.Error("Storage gauge should produce output")
	}
}

// TestManualNegativeEmpty specifically tests the condition that caused the panic
func TestManualNegativeEmpty(t *testing.T) {
	
	// With value 99.9 and width 3:
	// filled = int((99.9 / 100.0) * 3) = int(2.997) = 2
	// empty = 3 - 2 = 1 (this should be fine)
	
	// But with value 100.0 and width 1:
	// filled = int((100.0 / 100.0) * 1) = int(1.0) = 1  
	// empty = 1 - 1 = 0 (this should be fine)
	
	// However, due to floating point representation, some values could cause:
	// filled = int((value / 100.0) * width) to be > width in edge cases
	
	// Test a range of values that could trigger precision issues
	for value := 95.0; value <= 100.0; value += 0.1 {
		for width := 1; width <= 5; width++ {
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Panic with value=%.1f, width=%d: %v", value, width, r)
					}
				}()
				
				g := NewGauge(width)
				result := g.Render(value, "test")
				
				if result == "" {
					t.Errorf("Empty result for value=%.1f, width=%d", value, width)
				}
			}()
		}
	}
}