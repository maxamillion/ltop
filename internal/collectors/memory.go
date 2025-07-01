package collectors

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/admiller/ltop/internal/models"
	"github.com/admiller/ltop/internal/system"
)

type MemoryCollector struct {
	procReader *system.ProcReader
	sysReader  *system.SysReader
}

func NewMemoryCollector() *MemoryCollector {
	return &MemoryCollector{
		procReader: system.NewProcReader(),
		sysReader:  system.NewSysReader(),
	}
}

func (m *MemoryCollector) Collect() (*models.MemoryMetrics, error) {
	metrics := &models.MemoryMetrics{
		Timestamp: time.Now(),
		Details:   make(map[string]uint64),
	}

	memInfo, err := m.procReader.ReadMemInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to read memory info: %w", err)
	}

	if err := m.parseMemoryInfo(memInfo, metrics); err != nil {
		return nil, fmt.Errorf("failed to parse memory info: %w", err)
	}

	m.calculateMemoryUsage(metrics)
	m.calculateSwapUsage(metrics)

	return metrics, nil
}

func (m *MemoryCollector) parseMemoryInfo(memInfo map[string]string, metrics *models.MemoryMetrics) error {
	for key, value := range memInfo {
		valueKB, err := m.parseMemoryValueKB(value)
		if err != nil {
			continue
		}

		valueBytes := valueKB * 1024
		metrics.Details[key] = valueBytes

		switch key {
		case "MemTotal":
			metrics.Total = valueBytes
		case "MemFree":
			metrics.Free = valueBytes
		case "MemAvailable":
			metrics.Available = valueBytes
		case "Cached":
			metrics.Cached = valueBytes
		case "Buffers":
			metrics.Buffers = valueBytes
		case "Shmem":
			metrics.Shared = valueBytes
		case "SwapTotal":
			metrics.Swap.Total = valueBytes
		case "SwapFree":
			metrics.Swap.Free = valueBytes
		}
	}

	return nil
}

func (m *MemoryCollector) parseMemoryValueKB(value string) (uint64, error) {
	parts := strings.Fields(value)
	if len(parts) == 0 {
		return 0, fmt.Errorf("empty memory value")
	}

	numStr := parts[0]
	return strconv.ParseUint(numStr, 10, 64)
}

func (m *MemoryCollector) calculateMemoryUsage(metrics *models.MemoryMetrics) {
	if metrics.Available > 0 {
		metrics.Used = metrics.Total - metrics.Available
	} else {
		metrics.Used = metrics.Total - metrics.Free - metrics.Cached - metrics.Buffers
	}

	if metrics.Total > 0 {
		metrics.UsedPercent = float64(metrics.Used) / float64(metrics.Total) * 100.0
	}
}

func (m *MemoryCollector) calculateSwapUsage(metrics *models.MemoryMetrics) {
	metrics.Swap.Used = metrics.Swap.Total - metrics.Swap.Free

	if metrics.Swap.Total > 0 {
		metrics.Swap.UsedPercent = float64(metrics.Swap.Used) / float64(metrics.Swap.Total) * 100.0
	}
}
