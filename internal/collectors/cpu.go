package collectors

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/admiller/ltop/internal/models"
	"github.com/admiller/ltop/internal/system"
)

type CPUCollector struct {
	procReader   *system.ProcReader
	sysReader    *system.SysReader
	lastCPUTimes map[int]models.CPUTimes
	lastUpdate   time.Time
}

func NewCPUCollector() *CPUCollector {
	return &CPUCollector{
		procReader:   system.NewProcReader(),
		sysReader:    system.NewSysReader(),
		lastCPUTimes: make(map[int]models.CPUTimes),
		lastUpdate:   time.Now(),
	}
}

func (c *CPUCollector) Collect() (*models.CPUMetrics, error) {
	metrics := &models.CPUMetrics{
		Timestamp: time.Now(),
	}

	if err := c.collectCPUStats(metrics); err != nil {
		return nil, fmt.Errorf("failed to collect CPU stats: %w", err)
	}

	if err := c.collectLoadAverage(metrics); err != nil {
		return nil, fmt.Errorf("failed to collect load average: %w", err)
	}

	if err := c.collectCPUFrequency(metrics); err != nil {
		return metrics, nil
	}

	if err := c.collectCPUTemperature(metrics); err != nil {
		return metrics, nil
	}

	return metrics, nil
}

func (c *CPUCollector) collectCPUStats(metrics *models.CPUMetrics) error {
	lines, err := c.procReader.ReadCPUStat()
	if err != nil {
		return err
	}

	var totalTimes models.CPUTimes
	var cores []models.CPUCoreMetrics

	for _, line := range lines {
		if strings.HasPrefix(line, "cpu ") {
			times, err := c.parseCPULine(line)
			if err != nil {
				continue
			}
			totalTimes = times
			
			lastTimes, exists := c.lastCPUTimes[-1]
			if exists {
				metrics.Usage = c.calculateCPUUsage(lastTimes, times)
			}
			c.lastCPUTimes[-1] = times
		} else if strings.HasPrefix(line, "cpu") {
			cpuID, times, err := c.parseCPUCoreLine(line)
			if err != nil {
				continue
			}

			core := models.CPUCoreMetrics{
				ID:    cpuID,
				Times: times,
			}

			lastTimes, exists := c.lastCPUTimes[cpuID]
			if exists {
				core.Usage = c.calculateCPUUsage(lastTimes, times)
			}
			c.lastCPUTimes[cpuID] = times

			cores = append(cores, core)
		}
	}

	metrics.Times = totalTimes
	metrics.Cores = cores
	c.lastUpdate = time.Now()

	return nil
}

func (c *CPUCollector) parseCPULine(line string) (models.CPUTimes, error) {
	fields := strings.Fields(line)
	if len(fields) < 8 {
		return models.CPUTimes{}, fmt.Errorf("invalid CPU line format")
	}

	var times models.CPUTimes
	var err error

	if times.User, err = strconv.ParseUint(fields[1], 10, 64); err != nil {
		return times, err
	}
	if times.Nice, err = strconv.ParseUint(fields[2], 10, 64); err != nil {
		return times, err
	}
	if times.System, err = strconv.ParseUint(fields[3], 10, 64); err != nil {
		return times, err
	}
	if times.Idle, err = strconv.ParseUint(fields[4], 10, 64); err != nil {
		return times, err
	}
	if len(fields) > 5 {
		if times.IOWait, err = strconv.ParseUint(fields[5], 10, 64); err != nil {
			return times, err
		}
	}
	if len(fields) > 6 {
		if times.IRQ, err = strconv.ParseUint(fields[6], 10, 64); err != nil {
			return times, err
		}
	}
	if len(fields) > 7 {
		if times.SoftIRQ, err = strconv.ParseUint(fields[7], 10, 64); err != nil {
			return times, err
		}
	}
	if len(fields) > 8 {
		if times.Steal, err = strconv.ParseUint(fields[8], 10, 64); err != nil {
			return times, err
		}
	}
	if len(fields) > 9 {
		if times.Guest, err = strconv.ParseUint(fields[9], 10, 64); err != nil {
			return times, err
		}
	}
	if len(fields) > 10 {
		if times.GuestNice, err = strconv.ParseUint(fields[10], 10, 64); err != nil {
			return times, err
		}
	}

	return times, nil
}

func (c *CPUCollector) parseCPUCoreLine(line string) (int, models.CPUTimes, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return -1, models.CPUTimes{}, fmt.Errorf("invalid CPU core line format")
	}

	cpuStr := strings.TrimPrefix(fields[0], "cpu")
	cpuID, err := strconv.Atoi(cpuStr)
	if err != nil {
		return -1, models.CPUTimes{}, fmt.Errorf("invalid CPU ID: %s", cpuStr)
	}

	times, err := c.parseCPULine(line)
	return cpuID, times, err
}

func (c *CPUCollector) calculateCPUUsage(prev, curr models.CPUTimes) float64 {
	prevTotal := prev.User + prev.Nice + prev.System + prev.Idle + prev.IOWait + prev.IRQ + prev.SoftIRQ + prev.Steal
	currTotal := curr.User + curr.Nice + curr.System + curr.Idle + curr.IOWait + curr.IRQ + curr.SoftIRQ + curr.Steal

	prevIdle := prev.Idle + prev.IOWait
	currIdle := curr.Idle + curr.IOWait

	totalDiff := currTotal - prevTotal
	idleDiff := currIdle - prevIdle

	if totalDiff == 0 {
		return 0.0
	}

	return (float64(totalDiff-idleDiff) / float64(totalDiff)) * 100.0
}

func (c *CPUCollector) collectLoadAverage(metrics *models.CPUMetrics) error {
	loadStr, err := c.procReader.ReadLoadAvg()
	if err != nil {
		return err
	}

	fields := strings.Fields(loadStr)
	if len(fields) < 3 {
		return fmt.Errorf("invalid load average format")
	}

	for i := 0; i < 3; i++ {
		load, err := strconv.ParseFloat(fields[i], 64)
		if err != nil {
			return fmt.Errorf("failed to parse load average %d: %w", i, err)
		}
		metrics.LoadAverage[i] = load
	}

	return nil
}

func (c *CPUCollector) collectCPUFrequency(metrics *models.CPUMetrics) error {
	metrics.Frequency = make(map[string]uint64)

	for i, core := range metrics.Cores {
		freq, err := c.sysReader.ReadCPUFreq(strconv.Itoa(core.ID))
		if err != nil {
			continue
		}

		freqHz, err := strconv.ParseUint(freq, 10, 64)
		if err != nil {
			continue
		}

		cpuKey := fmt.Sprintf("cpu%d", i)
		metrics.Frequency[cpuKey] = freqHz * 1000
	}

	return nil
}

func (c *CPUCollector) collectCPUTemperature(metrics *models.CPUMetrics) error {
	tempStr, err := c.sysReader.ReadCPUTemperature()
	if err != nil {
		return err
	}

	tempMilliCelsius, err := strconv.ParseUint(tempStr, 10, 64)
	if err != nil {
		return err
	}

	metrics.Temperature = float64(tempMilliCelsius) / 1000.0
	return nil
}