package views

import (
	"fmt"
	"strings"

	"github.com/admiller/ltop/internal/models"
	"github.com/admiller/ltop/internal/ui/components"
	"github.com/admiller/ltop/internal/ui/styles"
	"github.com/admiller/ltop/pkg/utils"
)

type OverviewView struct {
	cpuGauge    *components.Gauge
	memoryGauge *components.Gauge
}

func NewOverviewView() *OverviewView {
	return &OverviewView{
		cpuGauge:    components.NewGauge(30),
		memoryGauge: components.NewGauge(30),
	}
}

func (ov *OverviewView) Render(snapshot *models.MetricsSnapshot, width, height int) string {
	if snapshot == nil {
		return "No data available"
	}

	var sections []string

	sections = append(sections, ov.renderSystemInfo(snapshot))
	sections = append(sections, ov.renderResourceSummary(snapshot))
	sections = append(sections, ov.renderTopProcesses(snapshot))

	content := strings.Join(sections, "\n\n")
	return styles.Panel().Width(width).Render(content)
}

func (ov *OverviewView) renderSystemInfo(snapshot *models.MetricsSnapshot) string {
	overview := snapshot.Overview
	
	var info []string
	info = append(info, styles.Title().Render("System Information"))
	
	if overview.Hostname != "" {
		info = append(info, fmt.Sprintf("Hostname: %s", styles.Info().Render(overview.Hostname)))
	}
	
	if overview.CurrentUser != "" {
		info = append(info, fmt.Sprintf("User: %s", styles.Info().Render(overview.CurrentUser)))
	}
	
	if overview.Uptime > 0 {
		info = append(info, fmt.Sprintf("Uptime: %s", styles.Info().Render(utils.FormatUptime(overview.Uptime))))
	}
	
	loadAvg := utils.FormatLoadAverage(snapshot.CPU.LoadAverage)
	info = append(info, fmt.Sprintf("Load Average: %s", styles.Info().Render(loadAvg)))
	
	info = append(info, fmt.Sprintf("Processes: %s (Running: %s, Sleeping: %s)",
		styles.Info().Render(fmt.Sprintf("%d", snapshot.Processes.Count)),
		styles.Success().Render(fmt.Sprintf("%d", snapshot.Processes.Running)),
		styles.Muted().Render(fmt.Sprintf("%d", snapshot.Processes.Sleeping))))

	return strings.Join(info, "\n")
}

func (ov *OverviewView) renderResourceSummary(snapshot *models.MetricsSnapshot) string {
	var summary []string
	summary = append(summary, styles.Title().Render("Resource Usage"))

	cpuGauge := ov.cpuGauge.Render(snapshot.CPU.Usage, "CPU")
	summary = append(summary, cpuGauge)

	memoryGauge := ov.memoryGauge.Render(snapshot.Memory.UsedPercent, "Memory")
	summary = append(summary, memoryGauge)
	
	memoryDetail := fmt.Sprintf("  %s / %s",
		utils.FormatBytes(snapshot.Memory.Used),
		utils.FormatBytes(snapshot.Memory.Total))
	summary = append(summary, styles.Muted().Render(memoryDetail))

	if snapshot.Memory.Swap.Total > 0 {
		swapGauge := ov.memoryGauge.Render(snapshot.Memory.Swap.UsedPercent, "Swap")
		summary = append(summary, swapGauge)
		
		swapDetail := fmt.Sprintf("  %s / %s",
			utils.FormatBytes(snapshot.Memory.Swap.Used),
			utils.FormatBytes(snapshot.Memory.Swap.Total))
		summary = append(summary, styles.Muted().Render(swapDetail))
	}

	if len(snapshot.Storage.Filesystems) > 0 {
		summary = append(summary, "")
		summary = append(summary, styles.Title().Render("Storage"))
		
		for _, fs := range snapshot.Storage.Filesystems {
			if fs.UsedPercent > 0 && len(fs.Mountpoint) > 1 {
				fsGauge := ov.memoryGauge.Render(fs.UsedPercent, utils.TruncateString(fs.Mountpoint, 15))
				summary = append(summary, fsGauge)
			}
		}
	}

	return strings.Join(summary, "\n")
}

func (ov *OverviewView) renderTopProcesses(snapshot *models.MetricsSnapshot) string {
	var processes []string
	processes = append(processes, styles.Title().Render("Top Processes"))

	if len(snapshot.Processes.Processes) == 0 {
		processes = append(processes, styles.Muted().Render("No processes found"))
		return strings.Join(processes, "\n")
	}

	headers := []string{"PID", "NAME", "CPU%", "MEMORY", "STATE"}
	processes = append(processes, ov.renderProcessHeader(headers))

	count := 10
	if len(snapshot.Processes.Processes) < count {
		count = len(snapshot.Processes.Processes)
	}

	for i := 0; i < count; i++ {
		proc := snapshot.Processes.Processes[i]
		row := ov.renderProcessRow(proc)
		processes = append(processes, row)
	}

	return strings.Join(processes, "\n")
}

func (ov *OverviewView) renderProcessHeader(headers []string) string {
	var parts []string
	widths := []int{8, 20, 8, 12, 8}
	
	for i, header := range headers {
		if i < len(widths) {
			text := utils.PadString(header, widths[i], ' ')
			parts = append(parts, styles.TableHeader().Render(text))
		}
	}
	
	return strings.Join(parts, " ")
}

func (ov *OverviewView) renderProcessRow(proc models.Process) string {
	var parts []string
	widths := []int{8, 20, 8, 12, 8}

	pid := utils.PadString(fmt.Sprintf("%d", proc.PID), widths[0], ' ')
	parts = append(parts, styles.TableRow().Render(pid))

	name := utils.TruncateString(proc.Name, widths[1])
	name = utils.PadString(name, widths[1], ' ')
	parts = append(parts, styles.TableRow().Render(name))

	cpu := utils.PadString(utils.FormatPercent(proc.CPUPercent), widths[2], ' ')
	cpuStyle := styles.PercentageColor(proc.CPUPercent)
	parts = append(parts, cpuStyle.Render(cpu))

	memory := utils.PadString(utils.FormatBytes(proc.MemoryRSS), widths[3], ' ')
	parts = append(parts, styles.TableRow().Render(memory))

	state := utils.PadString(utils.FormatProcessState(proc.State), widths[4], ' ')
	stateStyle := styles.ProcessStateColor(proc.State)
	parts = append(parts, stateStyle.Render(state))

	return strings.Join(parts, " ")
}