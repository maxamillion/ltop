package models

import (
	"testing"
	"time"
)

func TestDefaultSystemConfig(t *testing.T) {
	config := DefaultSystemConfig()
	
	if config.RefreshInterval <= 0 {
		t.Error("Refresh interval should be positive")
	}
	
	if config.MaxLogEntries <= 0 {
		t.Error("Max log entries should be positive")
	}
	
	if len(config.LogSources) == 0 {
		t.Error("Should have at least one log source")
	}
	
	if config.Theme == "" {
		t.Error("Theme should not be empty")
	}
}

func TestSystemConfigValidation(t *testing.T) {
	config := SystemConfig{
		RefreshInterval: time.Second,
		MaxLogEntries:   100,
		LogSources:      []string{"/var/log/syslog"},
		Theme:           "dark",
	}
	
	if config.RefreshInterval <= 0 {
		t.Error("Refresh interval should be positive")
	}
	
	if config.MaxLogEntries <= 0 {
		t.Error("Max log entries should be positive")
	}
	
}

func TestAppState(t *testing.T) {
	state := AppState{
		CurrentView: "overview",
		LastUpdate:  time.Now(),
		Paused:      false,
		Errors:      []AppError{},
	}
	
	if state.CurrentView == "" {
		t.Error("Current view should not be empty")
	}
	
	if state.LastUpdate.IsZero() {
		t.Error("Last update time should not be zero")
	}
	
	if state.Errors == nil {
		t.Error("Error slice should not be nil")
	}
}

func TestAppStateErrorHandling(t *testing.T) {
	state := AppState{
		Errors: []AppError{},
	}
	
	// Simulate adding errors
	errors := []AppError{
		{Message: "Failed to read CPU stats", Level: "error", Timestamp: time.Now()},
		{Message: "Network interface not found", Level: "warning", Timestamp: time.Now()},
		{Message: "Permission denied", Level: "error", Timestamp: time.Now()},
	}
	
	for _, err := range errors {
		state.Errors = append(state.Errors, err)
	}
	
	if len(state.Errors) != len(errors) {
		t.Errorf("Expected %d errors, got %d", len(errors), len(state.Errors))
	}
}

func TestSystemConfigThemeValidation(t *testing.T) {
	validThemes := []string{"dark", "light", "auto"}
	
	for _, theme := range validThemes {
		config := SystemConfig{
			Theme: theme,
		}
		
		if config.Theme != theme {
			t.Errorf("Theme should be %s, got %s", theme, config.Theme)
		}
	}
}

func TestSystemConfigRefreshInterval(t *testing.T) {
	testCases := []struct {
		interval time.Duration
		valid    bool
	}{
		{time.Millisecond * 100, true},   // 100ms - fast refresh
		{time.Second, true},              // 1s - normal
		{time.Second * 10, true},         // 10s - slow
		{time.Duration(0), false},        // 0 - invalid
		{time.Duration(-1), false},       // negative - invalid
	}
	
	for _, tc := range testCases {
		config := SystemConfig{
			RefreshInterval: tc.interval,
		}
		
		isValid := config.RefreshInterval > 0
		if isValid != tc.valid {
			t.Errorf("Refresh interval %v validity check failed. Expected: %v, Got: %v", 
				tc.interval, tc.valid, isValid)
		}
	}
}