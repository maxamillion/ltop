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

	// Always include overall CPU section (takes ~6 lines)
	overallSection := cv.renderOverallCPU(snapshot)
	sections = append(sections, overallSection)

	// Calculate remaining height for cores and info sections
	// Reserve space for panel padding and section separators
	usedHeight := strings.Count(overallSection, "\n") + 1 + 4 // +4 for padding and separators
	remainingHeight := height - usedHeight

	// Include CPU cores section if there's enough space (needs at least 4 lines)
	if remainingHeight >= 4 {
		coresSection := cv.renderCPUCores(snapshot, width)
		sections = append(sections, coresSection)
		
		coresHeight := strings.Count(coresSection, "\n") + 1 + 2 // +2 for separator
		remainingHeight -= coresHeight
		
		// Include CPU info section if there's still space (needs at least 6 lines)
		if remainingHeight >= 6 {
			sections = append(sections, cv.renderCPUInfo(snapshot))
		}
	}

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

func (cv *CPUView) renderCPUCores(snapshot *models.MetricsSnapshot, width int) string {
	var cores []string
	cores = append(cores, styles.Title().Render("Per-Core Usage"))

	if len(snapshot.CPU.Cores) == 0 {
		cores = append(cores, styles.Muted().Render("No core data available"))
		return strings.Join(cores, "\n")
	}

	// Calculate optimal number of columns based on terminal width
	// Each core entry needs about 50 characters (gauge + frequency)
	minColumnWidth := 50
	columnSpacing := 3 // Space between columns
	availableWidth := width - 4 // Account for panel padding
	
	numColumns := availableWidth / (minColumnWidth + columnSpacing)
	if numColumns == 0 {
		numColumns = 1
	}
	if numColumns > len(snapshot.CPU.Cores) {
		numColumns = len(snapshot.CPU.Cores)
	}

	// Calculate actual column width for uniform display
	actualColumnWidth := (availableWidth - (columnSpacing * (numColumns - 1))) / numColumns
	if actualColumnWidth < 45 {
		actualColumnWidth = 45 // Minimum readable width
	}

	// Ensure we have enough gauges for all cores
	for len(cv.coreGauges) < len(snapshot.CPU.Cores) {
		cv.coreGauges = append(cv.coreGauges, components.NewGauge(30))
	}

	// Prepare core data organized by rows and columns
	coresPerColumn := (len(snapshot.CPU.Cores) + numColumns - 1) / numColumns
	
	// Build core entries one by one, with frequency directly under each gauge
	var coreEntries []string
	
	for i, core := range snapshot.CPU.Cores {
		if i >= len(cv.coreGauges) {
			break
		}
		
		// Create gauge
		label := fmt.Sprintf("Core %d", core.ID)
		gaugeStr := cv.coreGauges[i].Render(core.Usage, label)
		
		// Create frequency string
		freqValue := ""
		if freq, exists := snapshot.CPU.Frequency[fmt.Sprintf("cpu%d", core.ID)]; exists {
			freqValue = utils.FormatHz(freq)
		} else {
			freqValue = "N/A"
		}
		freqStr := fmt.Sprintf("    %s", freqValue) // Simple indent
		
		// Combine gauge and frequency
		coreEntry := gaugeStr + "\n" + freqStr
		coreEntries = append(coreEntries, coreEntry)
	}
	
	// Arrange cores in columns - simplified approach
	var tableRows []string
	spacer := strings.Repeat(" ", columnSpacing)
	
	for row := 0; row < coresPerColumn; row++ {
		var gaugeRow []string
		var freqRow []string
		
		for col := 0; col < numColumns; col++ {
			coreIndex := col*coresPerColumn + row
			
			if coreIndex < len(coreEntries) {
				// Split the core entry back into gauge and freq lines
				lines := strings.Split(coreEntries[coreIndex], "\n")
				gaugeStr := utils.PadString(utils.TruncateString(lines[0], actualColumnWidth), actualColumnWidth, ' ')
				freqStr := utils.PadString(utils.TruncateString(lines[1], actualColumnWidth), actualColumnWidth, ' ')
				
				gaugeRow = append(gaugeRow, gaugeStr)
				freqRow = append(freqRow, freqStr)
			} else {
				// Empty column
				emptyCell := strings.Repeat(" ", actualColumnWidth)
				gaugeRow = append(gaugeRow, emptyCell)
				freqRow = append(freqRow, emptyCell)
			}
		}
		
		tableRows = append(tableRows, strings.Join(gaugeRow, spacer))
		tableRows = append(tableRows, strings.Join(freqRow, spacer))
		
		// Add spacing between core groups (except for the last row)
		if row < coresPerColumn-1 {
			tableRows = append(tableRows, "")
		}
	}
	
	cores = append(cores, strings.Join(tableRows, "\n"))
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
