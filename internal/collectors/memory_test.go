package collectors

import (
	"testing"
	"time"
)

func TestMemoryCollector(t *testing.T) {
	collector := NewMemoryCollector()
	if collector == nil {
		t.Fatal("NewMemoryCollector returned nil")
	}

	metrics, err := collector.Collect()
	if err != nil {
		t.Fatalf("Memory collection failed: %v", err)
	}

	if metrics == nil {
		t.Fatal("Memory metrics is nil")
	}

	if metrics.Total == 0 {
		t.Error("Total memory is zero")
	}

	if metrics.Used > metrics.Total {
		t.Errorf("Used memory (%d) exceeds total memory (%d)", metrics.Used, metrics.Total)
	}

	if metrics.UsedPercent < 0 || metrics.UsedPercent > 100 {
		t.Errorf("Invalid memory usage percentage: %f", metrics.UsedPercent)
	}

	if metrics.Timestamp.IsZero() {
		t.Error("Timestamp is zero")
	}

	if time.Since(metrics.Timestamp) > time.Minute {
		t.Error("Timestamp is too old")
	}

	if metrics.Details == nil {
		t.Error("Memory details map is nil")
	}
}

func TestMemoryCalculation(t *testing.T) {
	testCases := []struct {
		total    uint64
		used     uint64
		expected float64
	}{
		{1000, 500, 50.0},
		{2048, 1024, 50.0},
		{1000, 0, 0.0},
		{1000, 1000, 100.0},
	}

	for _, tc := range testCases {
		percentage := float64(tc.used) / float64(tc.total) * 100.0
		if abs(percentage-tc.expected) > 0.01 {
			t.Errorf("Memory percentage calculation incorrect. Total: %d, Used: %d, Expected: %f, Got: %f",
				tc.total, tc.used, tc.expected, percentage)
		}
	}
}
