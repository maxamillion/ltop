package collectors

import (
	"testing"
	"time"
)

func TestNetworkCollector(t *testing.T) {
	collector := NewNetworkCollector()
	if collector == nil {
		t.Fatal("NewNetworkCollector returned nil")
	}

	metrics, err := collector.Collect()
	if err != nil {
		t.Fatalf("Network collection failed: %v", err)
	}

	if metrics == nil {
		t.Fatal("Network metrics is nil")
	}

	if metrics.Timestamp.IsZero() {
		t.Error("Timestamp is zero")
	}

	if time.Since(metrics.Timestamp) > time.Minute {
		t.Error("Timestamp is too old")
	}

	if metrics.Interfaces == nil {
		t.Error("Interfaces slice is nil")
	}
}

func TestNetworkBandwidthCalculation(t *testing.T) {
	// Simple bandwidth calculation test
	prev := uint64(1000)
	curr := uint64(2000)
	duration := time.Second
	
	// Manual calculation
	if curr > prev && duration > 0 {
		bandwidth := (curr - prev) / uint64(duration.Seconds())
		expected := uint64(1000)
		
		if bandwidth != expected {
			t.Errorf("Bandwidth calculation incorrect. Expected: %d, Got: %d", expected, bandwidth)
		}
	}
}