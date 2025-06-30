package views

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/admiller/ltop/internal/models"
	"github.com/admiller/ltop/internal/ui/components"
	"github.com/admiller/ltop/internal/ui/styles"
	"github.com/admiller/ltop/pkg/utils"
)

type LogView struct {
	table        *components.Table
	filterLevel  string
	filterSource string
	autoScroll   bool
}

func NewLogView() *LogView {
	headers := []string{"TIME", "LEVEL", "SOURCE", "MESSAGE"}
	table := components.NewTable(headers)
	
	return &LogView{
		table:       table,
		filterLevel: "",
		autoScroll:  true,
	}
}

func (lv *LogView) Render(snapshot *models.MetricsSnapshot, width, height int) string {
	if snapshot == nil {
		return "No data available"
	}

	lv.table.SetSize(width-4, height-8)
	lv.table.ClearRows()

	filteredLogs := lv.filterLogs(snapshot.Logs.Entries)

	for _, entry := range filteredLogs {
		row := []string{
			utils.FormatTime(entry.Timestamp),
			strings.ToUpper(entry.Level),
			entry.Source,
			entry.Message,
		}
		lv.table.AddRow(row)
	}

	if lv.autoScroll && len(filteredLogs) > 0 {
		lv.table.SetSelected(len(filteredLogs) - 1)
	}

	var content strings.Builder
	
	content.WriteString(lv.renderLogSummary(snapshot.Logs))
	content.WriteString("\n\n")
	content.WriteString(lv.renderFilters())
	content.WriteString("\n")
	content.WriteString(lv.table.RenderWithInfo())

	return styles.Panel().Render(content.String())
}

func (lv *LogView) renderLogSummary(logs models.LogMetrics) string {
	var summary []string
	summary = append(summary, styles.Title().Render("System Logs"))

	errorStyle := styles.Error()
	warnStyle := styles.Warning()
	infoStyle := styles.Info()

	counts := fmt.Sprintf("Errors: %s | Warnings: %s | Info: %s",
		errorStyle.Render(strconv.Itoa(logs.ErrorCount)),
		warnStyle.Render(strconv.Itoa(logs.WarnCount)),
		infoStyle.Render(strconv.Itoa(logs.InfoCount)))
	summary = append(summary, counts)

	if len(logs.Sources) > 0 {
		sources := fmt.Sprintf("Sources: %s", strings.Join(logs.Sources, ", "))
		summary = append(summary, styles.Muted().Render(sources))
	}

	return strings.Join(summary, "\n")
}

func (lv *LogView) renderFilters() string {
	var filters []string
	
	if lv.filterLevel != "" {
		filters = append(filters, fmt.Sprintf("Level: %s", lv.filterLevel))
	}
	if lv.filterSource != "" {
		filters = append(filters, fmt.Sprintf("Source: %s", lv.filterSource))
	}
	
	scrollStatus := "Manual"
	if lv.autoScroll {
		scrollStatus = "Auto"
	}
	filters = append(filters, fmt.Sprintf("Scroll: %s", scrollStatus))

	if len(filters) > 0 {
		return styles.HelpText().Render("Filters: " + strings.Join(filters, " | "))
	}
	return styles.HelpText().Render("No filters active (f to set filters, a to toggle auto-scroll)")
}

func (lv *LogView) filterLogs(entries []models.LogEntry) []models.LogEntry {
	if lv.filterLevel == "" && lv.filterSource == "" {
		return entries
	}

	var filtered []models.LogEntry
	for _, entry := range entries {
		if lv.filterLevel != "" && !strings.EqualFold(entry.Level, lv.filterLevel) {
			continue
		}
		if lv.filterSource != "" && !strings.Contains(strings.ToLower(entry.Source), strings.ToLower(lv.filterSource)) {
			continue
		}
		filtered = append(filtered, entry)
	}
	return filtered
}

func (lv *LogView) MoveUp() {
	lv.autoScroll = false
	lv.table.MoveUp()
}

func (lv *LogView) MoveDown() {
	lv.autoScroll = false
	lv.table.MoveDown()
}

func (lv *LogView) PageUp() {
	lv.autoScroll = false
	lv.table.PageUp()
}

func (lv *LogView) PageDown() {
	lv.autoScroll = false
	lv.table.PageDown()
}

func (lv *LogView) ToggleAutoScroll() {
	lv.autoScroll = !lv.autoScroll
}

func (lv *LogView) SetLevelFilter(level string) {
	lv.filterLevel = level
}

func (lv *LogView) SetSourceFilter(source string) {
	lv.filterSource = source
}

func (lv *LogView) ClearFilters() {
	lv.filterLevel = ""
	lv.filterSource = ""
}

func (lv *LogView) GetSelectedEntry() []string {
	return lv.table.GetSelectedRow()
}