package models

import (
	"testing"
	"time"
)

func TestCPUMetrics(t *testing.T) {
	metrics := &CPUMetrics{
		Usage:       75.5,
		LoadAverage: [3]float64{1.0, 1.5, 2.0},
		Cores:       []CPUCoreMetrics{{ID: 0, Usage: 50.0}, {ID: 1, Usage: 60.0}},
		Frequency:   map[string]uint64{"cpu0": 2400000, "cpu1": 2400000},
		Temperature: 65.0,
		Timestamp:   time.Now(),
	}

	if metrics.Usage < 0 || metrics.Usage > 100 {
		t.Errorf("Invalid CPU usage: %f", metrics.Usage)
	}

	if len(metrics.LoadAverage) != 3 {
		t.Errorf("Expected 3 load average values, got %d", len(metrics.LoadAverage))
	}

	if len(metrics.Cores) <= 0 {
		t.Error("Core count should be positive")
	}

	if len(metrics.Frequency) <= 0 {
		t.Error("Frequency map should not be empty")
	}
}

func TestMemoryMetrics(t *testing.T) {
	metrics := &MemoryMetrics{
		Total:       8192,
		Used:        4096,
		Available:   4096,
		UsedPercent: 50.0,
		Swap: SwapMetrics{
			Total:       2048,
			Used:        512,
			UsedPercent: 25.0,
		},
		Timestamp: time.Now(),
	}

	if metrics.Used > metrics.Total {
		t.Error("Used memory cannot exceed total memory")
	}

	if metrics.Available > metrics.Total {
		t.Error("Available memory cannot exceed total memory")
	}

	if metrics.UsedPercent < 0 || metrics.UsedPercent > 100 {
		t.Errorf("Invalid memory usage percentage: %f", metrics.UsedPercent)
	}

	if metrics.Swap.Used > metrics.Swap.Total {
		t.Error("Used swap cannot exceed total swap")
	}
}

func TestStorageMetrics(t *testing.T) {
	fs := FilesystemMetrics{
		Mountpoint:  "/",
		FSType:      "ext4",
		Total:       1000000,
		Used:        500000,
		Free:        500000,
		UsedPercent: 50.0,
	}

	if fs.Used > fs.Total {
		t.Error("Used storage cannot exceed total storage")
	}

	if fs.Free > fs.Total {
		t.Error("Free storage cannot exceed total storage")
	}

	if fs.UsedPercent < 0 || fs.UsedPercent > 100 {
		t.Errorf("Invalid storage usage percentage: %f", fs.UsedPercent)
	}

	if fs.Mountpoint == "" {
		t.Error("Filesystem mountpoint cannot be empty")
	}
}

func TestNetworkMetrics(t *testing.T) {
	iface := NetworkInterface{
		Name:        "eth0",
		BytesRecv:   1000000,
		BytesSent:   500000,
		PacketsRecv: 10000,
		PacketsSent: 5000,
		ErrorsRecv:  0,
		ErrorsSent:  0,
		DroppedRecv: 0,
		DroppedSent: 0,
		Speed:       1000000000, // 1 Gbps
	}

	if iface.Name == "" {
		t.Error("Interface name cannot be empty")
	}

	if iface.BytesRecv == 0 || iface.BytesSent == 0 {
		t.Error("Bytes values cannot be zero")
	}

	if iface.PacketsRecv == 0 || iface.PacketsSent == 0 {
		t.Error("Packet values cannot be zero")
	}

	if iface.Speed <= 0 {
		t.Error("Interface speed must be positive")
	}
}

func TestProcessMetrics(t *testing.T) {
	proc := Process{
		PID:        1234,
		PPID:       1,
		Name:       "test-process",
		State:      "R",
		CPUPercent: 10.5,
		MemoryRSS:  1024000,
		MemoryVMS:  2048000,
		CPUTime:    time.Duration(5 * time.Minute),
		User:       "root",
		Command:    "/usr/bin/test-process",
		NumThreads: 4,
		Priority:   0,
		Nice:       0,
	}

	if proc.PID <= 0 {
		t.Error("PID should be positive")
	}

	if proc.PPID < 0 {
		t.Error("PPID cannot be negative")
	}

	if proc.Name == "" {
		t.Error("Process name cannot be empty")
	}

	if proc.CPUPercent < 0 || proc.CPUPercent > 100 {
		t.Errorf("Invalid CPU percentage: %f", proc.CPUPercent)
	}

	if proc.MemoryRSS == 0 {
		t.Error("Memory RSS cannot be zero")
	}

	if proc.NumThreads <= 0 {
		t.Error("Thread count should be positive")
	}
}

func TestLogMetrics(t *testing.T) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     "ERROR",
		Source:    "kernel",
		Message:   "Test error message",
		Service:   "ltop",
	}

	if entry.Timestamp.IsZero() {
		t.Error("Log timestamp cannot be zero")
	}

	if entry.Level == "" {
		t.Error("Log level cannot be empty")
	}

	if entry.Message == "" {
		t.Error("Log message cannot be empty")
	}
}

func TestMetricsSnapshot(t *testing.T) {
	snapshot := &MetricsSnapshot{
		Timestamp: time.Now(),
		CPU:       CPUMetrics{},
		Memory:    MemoryMetrics{},
		Storage:   StorageMetrics{},
		Network:   NetworkMetrics{},
		Processes: ProcessMetrics{},
		Logs:      LogMetrics{},
		Overview:  SystemOverview{},
	}

	if snapshot.Timestamp.IsZero() {
		t.Error("Snapshot timestamp cannot be zero")
	}
}
