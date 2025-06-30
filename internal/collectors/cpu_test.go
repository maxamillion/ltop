package collectors

import (
	"testing"
	"time"

	"github.com/admiller/ltop/internal/models"
)

func TestCPUCollector(t *testing.T) {
	collector := NewCPUCollector()
	if collector == nil {
		t.Fatal("NewCPUCollector returned nil")
	}

	metrics, err := collector.Collect()
	if err != nil {
		t.Fatalf("CPU collection failed: %v", err)
	}

	if metrics == nil {
		t.Fatal("CPU metrics is nil")
	}

	if metrics.Usage < 0 || metrics.Usage > 100 {
		t.Errorf("Invalid CPU usage: %f (should be 0-100)", metrics.Usage)
	}

	if len(metrics.LoadAverage) != 3 {
		t.Errorf("Expected 3 load average values, got %d", len(metrics.LoadAverage))
	}

	for i, load := range metrics.LoadAverage {
		if load < 0 {
			t.Errorf("Load average[%d] is negative: %f", i, load)
		}
	}

	if metrics.Timestamp.IsZero() {
		t.Error("Timestamp is zero")
	}

	if time.Since(metrics.Timestamp) > time.Minute {
		t.Error("Timestamp is too old")
	}
}

func TestCPUCalculation(t *testing.T) {
	collector := NewCPUCollector()
	
	prev := models.CPUTimes{
		User:   1000,
		Nice:   100,
		System: 500,
		Idle:   8000,
		IOWait: 400,
	}
	
	curr := models.CPUTimes{
		User:   1100,
		Nice:   100,
		System: 600,
		Idle:   8200,
		IOWait: 500,
	}
	
	usage := collector.calculateCPUUsage(prev, curr)
	
	if usage < 0 || usage > 100 {
		t.Errorf("Invalid CPU usage calculation: %f", usage)
	}
	
	// Total delta = 500, idle delta (idle + iowait) = 300, so non-idle = 200
	// Usage should be 200/500 = 40%
	expectedUsage := 40.0
	if abs(usage-expectedUsage) > 0.1 {
		t.Errorf("CPU usage calculation incorrect. Expected: %f, Got: %f", expectedUsage, usage)
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}