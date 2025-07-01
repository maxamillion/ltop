package views

import (
	"strings"

	"github.com/admiller/ltop/internal/models"
	"github.com/admiller/ltop/internal/ui/components"
	"github.com/admiller/ltop/internal/ui/styles"
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

func (v *LogView) Render(snapshot *models.MetricsSnapshot, width, height int) string {
	if snapshot == nil || len(snapshot.Logs.Entries) == 0 {
		return styles.Panel().Width(width).Height(height).Render("No log entries available")
	}

	v.updateEntries(snapshot.Logs.Entries)

	v.table.SetSize(width, height)
	content := v.table.Render()
	return styles.Panel().Width(width).Height(height).Render(content)
}

func (lv *LogView) updateEntries(entries []models.LogEntry) {
	filteredEntries := lv.filterLogs(entries)

	var rows [][]string
	for _, entry := range filteredEntries {
		rows = append(rows, []string{
			entry.Timestamp.Format("15:04:05"),
			strings.ToUpper(entry.Level),
			entry.Source,
			entry.Message,
		})
	}
	lv.table.Rows = rows

	if lv.autoScroll {
		lv.table.Selected = len(lv.table.Rows) - 1
		lv.table.ScrollTop = len(lv.table.Rows) - lv.table.Height + 1
		if lv.table.ScrollTop < 0 {
			lv.table.ScrollTop = 0
		}
	}
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
