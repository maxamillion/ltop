package views

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/admiller/ltop/internal/models"
	"github.com/admiller/ltop/internal/system"
	"github.com/admiller/ltop/internal/ui/components"
	"github.com/admiller/ltop/internal/ui/styles"
	"github.com/admiller/ltop/pkg/utils"
)

type ProcessView struct {
	table         *components.Table
	searchInput   *components.TextInput
	searchMode    bool
	sortField     string
	sortOrder     string
	processMgr    *system.ProcessManager
	confirmDialog *components.ConfirmDialog
	inputDialog   *components.InputDialog
	selectedPID   int
	actionType    string
}

func NewProcessView() *ProcessView {
	headers := []string{"PID", "PPID", "NAME", "STATE", "CPU%", "MEMORY", "TIME", "USER", "COMMAND"}
	table := components.NewTable(headers)
	searchInput := components.NewTextInput("Search processes...", 40)
	confirmDialog := components.NewConfirmDialog("Confirm Action", "")
	inputDialog := components.NewInputDialog("Process Management", "", "")
	
	return &ProcessView{
		table:         table,
		searchInput:   searchInput,
		searchMode:    false,
		sortField:     "cpu",
		sortOrder:     "desc",
		processMgr:    system.NewProcessManager(),
		confirmDialog: confirmDialog,
		inputDialog:   inputDialog,
	}
}

func (pv *ProcessView) Render(snapshot *models.MetricsSnapshot, width, height int) string {
	if snapshot == nil {
		return "No data available"
	}

	searchHeight := 0
	if pv.searchMode {
		searchHeight = 3
	}

	pv.table.SetSize(width-4, height-4-searchHeight)
	pv.table.ClearRows()

	filteredProcesses := pv.filterProcesses(snapshot.Processes.Processes)
	sortedProcesses := pv.sortProcesses(filteredProcesses)

	for _, proc := range sortedProcesses {
		row := []string{
			strconv.Itoa(proc.PID),
			strconv.Itoa(proc.PPID),
			proc.Name,
			utils.FormatProcessState(proc.State),
			utils.FormatPercent(proc.CPUPercent),
			utils.FormatBytes(proc.MemoryRSS),
			utils.FormatDuration(proc.CPUTime),
			proc.User,
			proc.Command,
		}
		pv.table.AddRow(row)
	}

	var content strings.Builder
	
	if pv.searchMode {
		content.WriteString("Search: ")
		content.WriteString(pv.searchInput.Render())
		content.WriteString("\n")
		content.WriteString(styles.HelpText().Render("Press Enter to search, Esc to cancel"))
		content.WriteString("\n\n")
	}

	content.WriteString(pv.renderProcessSummary(len(filteredProcesses), len(snapshot.Processes.Processes)))
	content.WriteString("\n")
	content.WriteString(pv.table.RenderWithInfo())

	result := styles.Panel().Render(content.String())

	if pv.confirmDialog.IsVisible() {
		result += "\n" + pv.confirmDialog.Render()
	}

	if pv.inputDialog.IsVisible() {
		result += "\n" + pv.inputDialog.Render()
	}

	return result
}

func (pv *ProcessView) MoveUp() {
	pv.table.MoveUp()
}

func (pv *ProcessView) MoveDown() {
	pv.table.MoveDown()
}

func (pv *ProcessView) PageUp() {
	pv.table.PageUp()
}

func (pv *ProcessView) PageDown() {
	pv.table.PageDown()
}

func (pv *ProcessView) GetSelectedProcess() []string {
	return pv.table.GetSelectedRow()
}

func (pv *ProcessView) renderProcessSummary(filtered, total int) string {
	summary := ""
	if filtered != total {
		summary = styles.Info().Render(
			fmt.Sprintf("Showing %d of %d processes", filtered, total))
	} else {
		summary = styles.Info().Render(
			fmt.Sprintf("Showing %d processes", total))
	}
	
	sortInfo := fmt.Sprintf("Sort: %s (%s)", pv.sortField, pv.sortOrder)
	return summary + " | " + styles.Muted().Render(sortInfo)
}

func (pv *ProcessView) filterProcesses(processes []models.Process) []models.Process {
	if pv.searchInput.IsEmpty() {
		return processes
	}

	query := strings.ToLower(pv.searchInput.GetValue())
	var filtered []models.Process

	for _, proc := range processes {
		if pv.processMatches(proc, query) {
			filtered = append(filtered, proc)
		}
	}

	return filtered
}

func (pv *ProcessView) processMatches(proc models.Process, query string) bool {
	fields := []string{
		proc.Name,
		proc.Command,
		proc.User,
		strconv.Itoa(proc.PID),
		utils.FormatProcessState(proc.State),
	}

	for _, field := range fields {
		if strings.Contains(strings.ToLower(field), query) {
			return true
		}
	}

	return false
}

func (pv *ProcessView) sortProcesses(processes []models.Process) []models.Process {
	if len(processes) == 0 {
		return processes
	}

	sorted := make([]models.Process, len(processes))
	copy(sorted, processes)

	switch pv.sortField {
	case "cpu":
		pv.sortByCPU(sorted)
	case "memory":
		pv.sortByMemory(sorted)
	case "pid":
		pv.sortByPID(sorted)
	case "name":
		pv.sortByName(sorted)
	case "time":
		pv.sortByTime(sorted)
	}

	if pv.sortOrder == "asc" {
		pv.reverseProcesses(sorted)
	}

	return sorted
}

func (pv *ProcessView) sortByCPU(processes []models.Process) {
	for i := 0; i < len(processes)-1; i++ {
		for j := i + 1; j < len(processes); j++ {
			if processes[i].CPUPercent < processes[j].CPUPercent {
				processes[i], processes[j] = processes[j], processes[i]
			}
		}
	}
}

func (pv *ProcessView) sortByMemory(processes []models.Process) {
	for i := 0; i < len(processes)-1; i++ {
		for j := i + 1; j < len(processes); j++ {
			if processes[i].MemoryRSS < processes[j].MemoryRSS {
				processes[i], processes[j] = processes[j], processes[i]
			}
		}
	}
}

func (pv *ProcessView) sortByPID(processes []models.Process) {
	for i := 0; i < len(processes)-1; i++ {
		for j := i + 1; j < len(processes); j++ {
			if processes[i].PID > processes[j].PID {
				processes[i], processes[j] = processes[j], processes[i]
			}
		}
	}
}

func (pv *ProcessView) sortByName(processes []models.Process) {
	for i := 0; i < len(processes)-1; i++ {
		for j := i + 1; j < len(processes); j++ {
			if strings.ToLower(processes[i].Name) > strings.ToLower(processes[j].Name) {
				processes[i], processes[j] = processes[j], processes[i]
			}
		}
	}
}

func (pv *ProcessView) sortByTime(processes []models.Process) {
	for i := 0; i < len(processes)-1; i++ {
		for j := i + 1; j < len(processes); j++ {
			if processes[i].CPUTime < processes[j].CPUTime {
				processes[i], processes[j] = processes[j], processes[i]
			}
		}
	}
}

func (pv *ProcessView) reverseProcesses(processes []models.Process) {
	for i := 0; i < len(processes)/2; i++ {
		j := len(processes) - 1 - i
		processes[i], processes[j] = processes[j], processes[i]
	}
}

func (pv *ProcessView) StartSearch() {
	pv.searchMode = true
	pv.searchInput.Focus()
	pv.searchInput.Clear()
}

func (pv *ProcessView) StopSearch() {
	pv.searchMode = false
	pv.searchInput.Blur()
}

func (pv *ProcessView) IsSearching() bool {
	return pv.searchMode
}

func (pv *ProcessView) HandleSearchInput(ch rune) {
	if pv.searchMode {
		pv.searchInput.InsertChar(ch)
	}
}

func (pv *ProcessView) HandleSearchBackspace() {
	if pv.searchMode {
		pv.searchInput.DeleteChar()
	}
}

func (pv *ProcessView) SetSortField(field string) {
	if pv.sortField == field {
		if pv.sortOrder == "desc" {
			pv.sortOrder = "asc"
		} else {
			pv.sortOrder = "desc"
		}
	} else {
		pv.sortField = field
		pv.sortOrder = "desc"
	}
}

func (pv *ProcessView) ToggleSortOrder() {
	if pv.sortOrder == "desc" {
		pv.sortOrder = "asc"
	} else {
		pv.sortOrder = "desc"
	}
}

func (pv *ProcessView) getSelectedPID() int {
	selectedRow := pv.table.GetSelectedRow()
	if len(selectedRow) > 0 {
		if pid, err := strconv.Atoi(selectedRow[0]); err == nil {
			return pid
		}
	}
	return 0
}

func (pv *ProcessView) ShowKillDialog() {
	pid := pv.getSelectedPID()
	if pid <= 0 {
		return
	}

	pv.selectedPID = pid
	pv.actionType = "kill"
	pv.confirmDialog.Title = "Kill Process"
	pv.confirmDialog.Message = fmt.Sprintf("Are you sure you want to terminate process %d?", pid)
	pv.confirmDialog.Show()
}

func (pv *ProcessView) ShowForceKillDialog() {
	pid := pv.getSelectedPID()
	if pid <= 0 {
		return
	}

	pv.selectedPID = pid
	pv.actionType = "force_kill"
	pv.confirmDialog.Title = "Force Kill Process"
	pv.confirmDialog.Message = fmt.Sprintf("Are you sure you want to force kill process %d? This cannot be undone.", pid)
	pv.confirmDialog.Show()
}

func (pv *ProcessView) ShowStopDialog() {
	pid := pv.getSelectedPID()
	if pid <= 0 {
		return
	}

	pv.selectedPID = pid
	pv.actionType = "stop"
	pv.confirmDialog.Title = "Stop Process"
	pv.confirmDialog.Message = fmt.Sprintf("Are you sure you want to stop process %d?", pid)
	pv.confirmDialog.Show()
}

func (pv *ProcessView) ShowContinueDialog() {
	pid := pv.getSelectedPID()
	if pid <= 0 {
		return
	}

	pv.selectedPID = pid
	pv.actionType = "continue"
	pv.confirmDialog.Title = "Continue Process"
	pv.confirmDialog.Message = fmt.Sprintf("Continue process %d?", pid)
	pv.confirmDialog.Show()
}

func (pv *ProcessView) ShowNiceDialog() {
	pid := pv.getSelectedPID()
	if pid <= 0 {
		return
	}

	pv.selectedPID = pid
	pv.actionType = "nice"
	pv.inputDialog.Title = "Change Process Priority"
	pv.inputDialog.Message = fmt.Sprintf("Enter new priority for process %d\n(Range: -20 to 19, lower = higher priority)", pid)
	pv.inputDialog.Show()
}

func (pv *ProcessView) ExecuteAction() error {
	switch pv.actionType {
	case "kill":
		return pv.processMgr.TerminateProcess(pv.selectedPID)
	case "force_kill":
		return pv.processMgr.ForceKillProcess(pv.selectedPID)
	case "stop":
		return pv.processMgr.StopProcess(pv.selectedPID)
	case "continue":
		return pv.processMgr.ContinueProcess(pv.selectedPID)
	case "nice":
		if priority, err := strconv.Atoi(pv.inputDialog.GetValue()); err == nil {
			if priority >= -20 && priority <= 19 {
				return pv.processMgr.SetProcessPriority(pv.selectedPID, priority)
			}
			return fmt.Errorf("priority must be between -20 and 19")
		}
		return fmt.Errorf("invalid priority value")
	}
	return fmt.Errorf("unknown action: %s", pv.actionType)
}

func (pv *ProcessView) IsDialogActive() bool {
	return pv.confirmDialog.IsVisible() || pv.inputDialog.IsVisible()
}

func (pv *ProcessView) HandleDialogInput(key string) error {
	if pv.confirmDialog.IsVisible() {
		switch key {
		case "left", "h":
			pv.confirmDialog.MoveLeft()
		case "right", "l":
			pv.confirmDialog.MoveRight()
		case "enter":
			if pv.confirmDialog.IsConfirmSelected() {
				err := pv.ExecuteAction()
				pv.confirmDialog.Hide()
				return err
			}
			pv.confirmDialog.Hide()
		case "esc":
			pv.confirmDialog.Hide()
		}
	} else if pv.inputDialog.IsVisible() {
		switch key {
		case "enter":
			err := pv.ExecuteAction()
			pv.inputDialog.Hide()
			return err
		case "esc":
			pv.inputDialog.Hide()
		case "backspace":
			pv.inputDialog.HandleBackspace()
		default:
			if len(key) == 1 {
				pv.inputDialog.HandleInput(rune(key[0]))
			}
		}
	}
	return nil
}