package views

import (
	"fmt"
	"strings"

	"github.com/admiller/ltop/internal/models"
	"github.com/admiller/ltop/internal/ui/components"
	"github.com/admiller/ltop/internal/ui/styles"
	"github.com/admiller/ltop/pkg/utils"
)

type CPUView struct {
	overallGauge *components.Gauge
	coreGauges   []*components.Gauge
}

func NewCPUView() *CPUView {
	return &CPUView{
		overallGauge: components.NewGauge(40),
		coreGauges:   make([]*components.Gauge, 0),
	}
}

func (cv *CPUView) Render(snapshot *models.MetricsSnapshot, width, height int) string {
	if snapshot == nil {
		return "No data available"
	}

	var sections []string

	sections = append(sections, cv.renderOverallCPU(snapshot))
	sections = append(sections, cv.renderCPUCores(snapshot))
	sections = append(sections, cv.renderCPUInfo(snapshot))

	content := strings.Join(sections, "\n\n")
	return styles.Panel().Width(width).Height(height).Render(content)
}

func (cv *CPUView) renderOverallCPU(snapshot *models.MetricsSnapshot) string {
	var info []string
	info = append(info, styles.Title().Render("CPU Usage"))

	overallGauge := cv.overallGauge.Render(snapshot.CPU.Usage, "Overall")
	info = append(info, overallGauge)

	loadAvg := utils.FormatLoadAverage(snapshot.CPU.LoadAverage)
	info = append(info, fmt.Sprintf("Load Average: %s", styles.Info().Render(loadAvg)))

	if snapshot.CPU.Temperature > 0 {
		temp := utils.FormatTemperature(snapshot.CPU.Temperature)
		tempStyle := styles.Info()
		if snapshot.CPU.Temperature > 80 {
			tempStyle = styles.Error()
		} else if snapshot.CPU.Temperature > 70 {
			tempStyle = styles.Warning()
		}
		info = append(info, fmt.Sprintf("Temperature: %s", tempStyle.Render(temp)))
	}

	return strings.Join(info, "\n")
}

func (cv *CPUView) renderCPUCores(snapshot *models.MetricsSnapshot) string {
	var cores []string
	cores = append(cores, styles.Title().Render("Per-Core Usage"))

	for len(cv.coreGauges) < len(snapshot.CPU.Cores) {
		cv.coreGauges = append(cv.coreGauges, components.NewGauge(30))
	}

	for i, core := range snapshot.CPU.Cores {
		if i < len(cv.coreGauges) {
			label := fmt.Sprintf("Core %d", core.ID)
			gauge := cv.coreGauges[i].Render(core.Usage, label)
			cores = append(cores, gauge)

			if freq, exists := snapshot.CPU.Frequency[fmt.Sprintf("cpu%d", i)]; exists {
				freqStr := utils.FormatHz(freq)
				cores = append(cores, fmt.Sprintf("  Frequency: %s", styles.Muted().Render(freqStr)))
			}
		}
	}

	return strings.Join(cores, "\n")
}

func (cv *CPUView) renderCPUInfo(snapshot *models.MetricsSnapshot) string {
	var info []string
	info = append(info, styles.Title().Render("CPU Statistics"))

	times := snapshot.CPU.Times
	total := times.User + times.Nice + times.System + times.Idle + times.IOWait + times.IRQ + times.SoftIRQ + times.Steal

	if total > 0 {
		userPct := float64(times.User) / float64(total) * 100
		systemPct := float64(times.System) / float64(total) * 100
		iowaitPct := float64(times.IOWait) / float64(total) * 100

		info = append(info, fmt.Sprintf("User: %s", utils.FormatPercent(userPct)))
		info = append(info, fmt.Sprintf("System: %s", utils.FormatPercent(systemPct)))
		info = append(info, fmt.Sprintf("I/O Wait: %s", utils.FormatPercent(iowaitPct)))
	}

	return strings.Join(info, "\n")
}
