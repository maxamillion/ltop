package collectors

import (
	"testing"
	"time"
)

func TestStorageCollector(t *testing.T) {
	collector := NewStorageCollector()
	if collector == nil {
		t.Fatal("NewStorageCollector returned nil")
	}

	metrics, err := collector.Collect()
	if err != nil {
		t.Fatalf("Storage collection failed: %v", err)
	}

	if metrics == nil {
		t.Fatal("Storage metrics is nil")
	}

	if metrics.Timestamp.IsZero() {
		t.Error("Timestamp is zero")
	}

	if time.Since(metrics.Timestamp) > time.Minute {
		t.Error("Timestamp is too old")
	}

	if metrics.Filesystems == nil {
		t.Error("Filesystems slice is nil")
	}

	if metrics.IOStats == nil {
		t.Error("IOStats map is nil")
	}
}

func TestStorageCalculation(t *testing.T) {
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
			t.Errorf("Storage percentage calculation incorrect. Total: %d, Used: %d, Expected: %f, Got: %f",
				tc.total, tc.used, tc.expected, percentage)
		}
	}
}
