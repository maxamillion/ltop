package app

import (
	"testing"
	"time"
)

func TestNewApp(t *testing.T) {
	app := New()
	if app == nil {
		t.Fatal("New() returned nil")
	}

	if app.cpuCollector == nil {
		t.Error("CPU collector is nil")
	}

	if app.memoryCollector == nil {
		t.Error("Memory collector is nil")
	}

	if app.processCollector == nil {
		t.Error("Process collector is nil")
	}

	if app.storageCollector == nil {
		t.Error("Storage collector is nil")
	}

	if app.networkCollector == nil {
		t.Error("Network collector is nil")
	}

	if app.logCollector == nil {
		t.Error("Log collector is nil")
	}

	if app.ctx == nil {
		t.Error("Context is nil")
	}

	if app.cancel == nil {
		t.Error("Cancel function is nil")
	}
}

func TestAppCollectMetrics(t *testing.T) {
	app := New()

	err := app.CollectMetrics()
	if err != nil {
		t.Fatalf("CollectMetrics failed: %v", err)
	}

	snapshot := app.GetLastSnapshot()

	if snapshot == nil {
		t.Fatal("Snapshot is nil")
	}

	// Metrics are value types, not pointers, so just check timestamp
	if snapshot.CPU.Timestamp.IsZero() && snapshot.Memory.Timestamp.IsZero() {
		t.Error("No metrics were collected")
	}

	if snapshot.Timestamp.IsZero() {
		t.Error("Snapshot timestamp is zero")
	}

	if time.Since(snapshot.Timestamp) > time.Minute {
		t.Error("Snapshot timestamp is too old")
	}
}

func TestAppGetLastSnapshot(t *testing.T) {
	app := New()

	// Initially should be nil
	snapshot := app.GetLastSnapshot()
	if snapshot != nil {
		t.Error("Initial snapshot should be nil")
	}

	// After collecting metrics, should not be nil
	err := app.CollectMetrics()
	if err != nil {
		t.Fatalf("CollectMetrics failed: %v", err)
	}

	snapshot = app.GetLastSnapshot()
	if snapshot == nil {
		t.Error("Snapshot should not be nil after collection")
	}
}

func TestAppSetConfig(t *testing.T) {
	app := New()

	originalInterval := app.config.RefreshInterval
	newInterval := time.Second * 5

	newConfig := app.config
	newConfig.RefreshInterval = newInterval
	app.SetConfig(newConfig)

	if app.config.RefreshInterval != newInterval {
		t.Errorf("Config update failed. Expected: %v, Got: %v", newInterval, app.config.RefreshInterval)
	}

	// Restore original
	newConfig.RefreshInterval = originalInterval
	app.SetConfig(newConfig)
}

func TestAppGetConfig(t *testing.T) {
	app := New()

	config := app.GetConfig()
	if config.RefreshInterval <= 0 {
		t.Error("Refresh interval should be positive")
	}

	if config.MaxLogEntries <= 0 {
		t.Error("Max log entries should be positive")
	}

	if len(config.LogSources) == 0 {
		t.Error("Should have at least one log source")
	}
}

func TestAppGetState(t *testing.T) {
	app := New()

	state := app.GetState()

	// Initially, state should have default values
	if state.CurrentView != "" {
		t.Error("Current view should be empty")
	}

	// Errors slice might be nil initially, that's OK
	if len(state.Errors) != 0 {
		t.Error("Error slice should be empty")
	}
}

func TestAppSetCurrentView(t *testing.T) {
	app := New()

	testViews := []string{"overview", "cpu", "memory", "storage", "network", "processes", "logs"}

	for _, view := range testViews {
		app.SetCurrentView(view)

		state := app.GetState()
		if state.CurrentView != view {
			t.Errorf("Expected current view to be %s, got %s", view, state.CurrentView)
		}
	}
}

func TestAppTogglePause(t *testing.T) {
	app := New()

	// Initially should not be paused
	state := app.GetState()
	initialPaused := state.Paused

	// Toggle pause
	app.TogglePause()

	state = app.GetState()
	if state.Paused == initialPaused {
		t.Error("Pause state should have toggled")
	}

	// Toggle back
	app.TogglePause()

	state = app.GetState()
	if state.Paused != initialPaused {
		t.Error("Pause state should have toggled back")
	}
}

func TestAppContext(t *testing.T) {
	app := New()

	// Context should not be cancelled initially
	select {
	case <-app.ctx.Done():
		t.Error("Context should not be cancelled initially")
	default:
		// OK
	}

	// Cancel the context
	app.cancel()

	// Context should now be cancelled
	select {
	case <-app.ctx.Done():
		// OK
	case <-time.After(time.Millisecond * 100):
		t.Error("Context should be cancelled after calling cancel")
	}
}
