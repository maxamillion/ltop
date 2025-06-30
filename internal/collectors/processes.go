package collectors

import (
	"fmt"
	"os/user"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/admiller/ltop/internal/models"
	"github.com/admiller/ltop/internal/system"
)

type ProcessCollector struct {
	procReader   *system.ProcReader
	lastCPUTimes map[int]uint64
	lastUpdate   time.Time
	totalCPU     uint64
	lastTotalCPU uint64
}

func NewProcessCollector() *ProcessCollector {
	return &ProcessCollector{
		procReader:   system.NewProcReader(),
		lastCPUTimes: make(map[int]uint64),
		lastUpdate:   time.Now(),
	}
}

func (p *ProcessCollector) Collect() (*models.ProcessMetrics, error) {
	metrics := &models.ProcessMetrics{
		Timestamp: time.Now(),
		Processes: make([]models.Process, 0),
	}

	pids, err := p.procReader.ReadProcesses()
	if err != nil {
		return nil, fmt.Errorf("failed to read processes: %w", err)
	}

	totalCPU, err := p.getTotalCPUTime()
	if err != nil {
		return nil, fmt.Errorf("failed to get total CPU time: %w", err)
	}

	processes := make([]models.Process, 0, len(pids))
	stateCounts := make(map[string]int)

	for _, pid := range pids {
		process, err := p.collectProcess(pid, totalCPU)
		if err != nil {
			continue
		}

		processes = append(processes, process)
		stateCounts[process.State]++
	}

	sort.Slice(processes, func(i, j int) bool {
		return processes[i].CPUPercent > processes[j].CPUPercent
	})

	metrics.Processes = processes
	metrics.Count = len(processes)
	metrics.Running = stateCounts["R"]
	metrics.Sleeping = stateCounts["S"] + stateCounts["I"]
	metrics.Stopped = stateCounts["T"]
	metrics.Zombie = stateCounts["Z"]

	p.lastTotalCPU = p.totalCPU
	p.totalCPU = totalCPU
	p.lastUpdate = time.Now()

	return metrics, nil
}

func (p *ProcessCollector) collectProcess(pid string, totalCPU uint64) (models.Process, error) {
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		return models.Process{}, err
	}

	process := models.Process{
		PID: pidInt,
	}

	if err := p.collectProcessStat(pid, &process); err != nil {
		return process, err
	}

	if err := p.collectProcessStatus(pid, &process); err != nil {
		return process, nil
	}

	if err := p.collectProcessCmdline(pid, &process); err != nil {
		return process, nil
	}

	p.calculateCPUPercent(&process, totalCPU)

	return process, nil
}

func (p *ProcessCollector) collectProcessStat(pid string, process *models.Process) error {
	statLine, err := p.procReader.ReadProcessStat(pid)
	if err != nil {
		return err
	}

	fields := strings.Fields(statLine)
	if len(fields) < 52 {
		return fmt.Errorf("invalid stat format for PID %s", pid)
	}

	if ppid, err := strconv.Atoi(fields[3]); err == nil {
		process.PPID = ppid
	}

	process.State = fields[2]
	process.Name = strings.Trim(fields[1], "()")

	if priority, err := strconv.Atoi(fields[17]); err == nil {
		process.Priority = priority
	}

	if nice, err := strconv.Atoi(fields[18]); err == nil {
		process.Nice = nice
	}

	if numThreads, err := strconv.Atoi(fields[19]); err == nil {
		process.NumThreads = numThreads
	}

	if utime, err := strconv.ParseUint(fields[13], 10, 64); err == nil {
		if stime, err := strconv.ParseUint(fields[14], 10, 64); err == nil {
			totalTime := utime + stime
			process.CPUTime = time.Duration(totalTime) * time.Millisecond / 100
		}
	}

	if rss, err := strconv.ParseUint(fields[23], 10, 64); err == nil {
		process.MemoryRSS = rss * 4096
	}

	if vsize, err := strconv.ParseUint(fields[22], 10, 64); err == nil {
		process.MemoryVMS = vsize
	}

	if startTime, err := strconv.ParseUint(fields[21], 10, 64); err == nil {
		bootTime := p.getBootTime()
		process.CreateTime = bootTime.Add(time.Duration(startTime) * time.Second / 100)
	}

	return nil
}

func (p *ProcessCollector) collectProcessStatus(pid string, process *models.Process) error {
	status, err := p.procReader.ReadProcessStatus(pid)
	if err != nil {
		return err
	}

	if uid, exists := status["Uid"]; exists {
		fields := strings.Fields(uid)
		if len(fields) > 0 {
			if uidInt, err := strconv.Atoi(fields[0]); err == nil {
				if u, err := user.LookupId(strconv.Itoa(uidInt)); err == nil {
					process.User = u.Username
				}
			}
		}
	}

	if gid, exists := status["Gid"]; exists {
		fields := strings.Fields(gid)
		if len(fields) > 0 {
			if gidInt, err := strconv.Atoi(fields[0]); err == nil {
				if g, err := user.LookupGroupId(strconv.Itoa(gidInt)); err == nil {
					process.Group = g.Name
				}
			}
		}
	}

	if fdSize, exists := status["FDSize"]; exists {
		if fds, err := strconv.Atoi(strings.Fields(fdSize)[0]); err == nil {
			process.NumFDs = fds
		}
	}

	return nil
}

func (p *ProcessCollector) collectProcessCmdline(pid string, process *models.Process) error {
	cmdline, err := p.procReader.ReadProcessCmdline(pid)
	if err != nil {
		return err
	}

	if cmdline == "" {
		process.Command = process.Name
	} else {
		process.Command = cmdline
	}

	return nil
}

func (p *ProcessCollector) calculateCPUPercent(process *models.Process, totalCPU uint64) {
	cpuTime := uint64(process.CPUTime / time.Millisecond)
	
	lastCPUTime, exists := p.lastCPUTimes[process.PID]
	if !exists {
		p.lastCPUTimes[process.PID] = cpuTime
		return
	}

	cpuDelta := cpuTime - lastCPUTime
	totalDelta := totalCPU - p.lastTotalCPU

	if totalDelta > 0 {
		process.CPUPercent = float64(cpuDelta) / float64(totalDelta) * 100.0
	}

	p.lastCPUTimes[process.PID] = cpuTime
}

func (p *ProcessCollector) getTotalCPUTime() (uint64, error) {
	lines, err := p.procReader.ReadCPUStat()
	if err != nil {
		return 0, err
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)
			if len(fields) < 8 {
				continue
			}

			var total uint64
			for i := 1; i < 8; i++ {
				if val, err := strconv.ParseUint(fields[i], 10, 64); err == nil {
					total += val
				}
			}
			return total, nil
		}
	}

	return 0, fmt.Errorf("CPU total time not found")
}

func (p *ProcessCollector) getBootTime() time.Time {
	uptime, err := p.procReader.ReadUptime()
	if err != nil {
		return time.Now()
	}

	fields := strings.Fields(uptime)
	if len(fields) == 0 {
		return time.Now()
	}

	uptimeSeconds, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return time.Now()
	}

	return time.Now().Add(-time.Duration(uptimeSeconds) * time.Second)
}