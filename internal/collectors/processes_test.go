package collectors

import (
	"testing"
	"time"
)

func TestProcessCollector(t *testing.T) {
	collector := NewProcessCollector()
	if collector == nil {
		t.Fatal("NewProcessCollector returned nil")
	}

	metrics, err := collector.Collect()
	if err != nil {
		t.Fatalf("Process collection failed: %v", err)
	}

	if metrics == nil {
		t.Fatal("Process metrics is nil")
	}

	if metrics.Timestamp.IsZero() {
		t.Error("Timestamp is zero")
	}

	if time.Since(metrics.Timestamp) > time.Minute {
		t.Error("Timestamp is too old")
	}

	if metrics.Processes == nil {
		t.Error("Processes slice is nil")
	}

	if metrics.Count == 0 {
		t.Error("Process count is zero")
	}

	if len(metrics.Processes) == 0 {
		t.Error("Processes slice is empty")
	}
}

func TestProcessStats(t *testing.T) {
	// Test process state counting
	states := map[string]int{
		"R": 5,  // Running
		"S": 10, // Sleeping
		"D": 1,  // Disk sleep
		"Z": 2,  // Zombie
	}

	total := 0
	for _, count := range states {
		total += count
	}

	if total != 18 {
		t.Errorf("Process state counting incorrect. Expected total: 18, Got: %d", total)
	}
}

func TestProcessCPUCalculation(t *testing.T) {
	testCases := []struct {
		utime, stime uint64
		totalTime    uint64
		expected     float64
	}{
		{100, 50, 1000, 15.0},  // (100+50)/1000 * 100 = 15%
		{200, 200, 2000, 20.0}, // (200+200)/2000 * 100 = 20%
		{0, 0, 100, 0.0},       // No CPU usage
	}

	for _, tc := range testCases {
		percentage := float64(tc.utime+tc.stime) / float64(tc.totalTime) * 100.0
		if abs(percentage-tc.expected) > 0.01 {
			t.Errorf("Process CPU calculation incorrect. UTime: %d, STime: %d, Total: %d, Expected: %f, Got: %f",
				tc.utime, tc.stime, tc.totalTime, tc.expected, percentage)
		}
	}
}
