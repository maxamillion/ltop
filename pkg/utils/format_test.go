package utils

import (
	"testing"
	"time"
)

func TestFormatBytes(t *testing.T) {
	testCases := []struct {
		input    uint64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
		{1099511627776, "1.0 TB"},
	}

	for _, tc := range testCases {
		result := FormatBytes(tc.input)
		if result != tc.expected {
			t.Errorf("FormatBytes(%d) = %s; expected %s", tc.input, result, tc.expected)
		}
	}
}

func TestFormatPercent(t *testing.T) {
	testCases := []struct {
		input    float64
		expected string
	}{
		{0.0, "0.0%"},
		{50.0, "50.0%"},
		{100.0, "100.0%"},
		{33.33, "33.3%"},
		{99.99, "100.0%"},
	}

	for _, tc := range testCases {
		result := FormatPercent(tc.input)
		if result != tc.expected {
			t.Errorf("FormatPercent(%f) = %s; expected %s", tc.input, result, tc.expected)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	testCases := []struct {
		input    time.Duration
		expected string
	}{
		{30 * time.Second, "30.0s"},
		{90 * time.Second, "1.5m"},
		{3600 * time.Second, "1.0h"},
		{25*time.Hour + 30*time.Minute, "1d 1h"},
	}

	for _, tc := range testCases {
		result := FormatDuration(tc.input)
		if result != tc.expected {
			t.Errorf("FormatDuration(%v) = %s; expected %s", tc.input, result, tc.expected)
		}
	}
}

func TestTruncateString(t *testing.T) {
	testCases := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "he..."},
		{"test", 3, "tes"},
		{"a", 1, "a"},
		{"", 5, ""},
	}

	for _, tc := range testCases {
		result := TruncateString(tc.input, tc.maxLen)
		if result != tc.expected {
			t.Errorf("TruncateString(%s, %d) = %s; expected %s", tc.input, tc.maxLen, result, tc.expected)
		}
	}
}

func TestFormatHz(t *testing.T) {
	testCases := []struct {
		input    uint64
		expected string
	}{
		{500, "500 Hz"},
		{1500, "1.5 KHz"},
		{2000000, "2.0 MHz"},
		{3000000000, "3.0 GHz"},
	}

	for _, tc := range testCases {
		result := FormatHz(tc.input)
		if result != tc.expected {
			t.Errorf("FormatHz(%d) = %s; expected %s", tc.input, result, tc.expected)
		}
	}
}

func TestFormatProcessState(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"R", "Running"},
		{"S", "Sleeping"},
		{"D", "Disk Sleep"},
		{"T", "Stopped"},
		{"Z", "Zombie"},
		{"X", "Dead"},
	}

	for _, tc := range testCases {
		result := FormatProcessState(tc.input)
		if result != tc.expected {
			t.Errorf("FormatProcessState(%s) = %s; expected %s", tc.input, result, tc.expected)
		}
	}
}