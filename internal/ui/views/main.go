package views

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/admiller/ltop/internal/app"
	"github.com/admiller/ltop/internal/models"
	"github.com/admiller/ltop/internal/ui/styles"
	"github.com/admiller/ltop/pkg/utils"
)

type Model struct {
	app          *app.App
	currentView  models.ViewType
	width        int
	height       int
	overviewView *OverviewView
	cpuView      *CPUView
	memoryView   *MemoryView
	storageView  *StorageView
	networkView  *NetworkView
	processView  *ProcessView
	logView      *LogView
	lastUpdate   time.Time
	showHelp     bool
	err          error
}

type TickMsg time.Time

func NewModel(ltopApp *app.App) Model {
	return Model{
		app:          ltopApp,
		currentView:  models.ViewOverview,
		overviewView: NewOverviewView(),
		cpuView:      NewCPUView(),
		memoryView:   NewMemoryView(),
		storageView:  NewStorageView(),
		networkView:  NewNetworkView(),
		processView:  NewProcessView(),
		logView:      NewLogView(),
		lastUpdate:   time.Now(),
		showHelp:     false,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		tea.EnterAltScreen,
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "h", "?":
			m.showHelp = !m.showHelp
			return m, nil
		case "p":
			m.app.TogglePause()
			return m, nil
		case "r":
			config := m.app.GetConfig()
			if config.RefreshInterval == time.Second {
				config.RefreshInterval = 5 * time.Second
			} else {
				config.RefreshInterval = time.Second
			}
			m.app.SetConfig(config)
			return m, nil
		case "T":
			config := m.app.GetConfig()
			if config.Theme == "dark" {
				config.Theme = "light"
			} else {
				config.Theme = "dark"
			}
			m.app.SetConfig(config)
			return m, nil
		case "1":
			m.currentView = models.ViewOverview
			return m, nil
		case "2":
			m.currentView = models.ViewCPU
			return m, nil
		case "3":
			m.currentView = models.ViewMemory
			return m, nil
		case "4":
			m.currentView = models.ViewStorage
			return m, nil
		case "5":
			m.currentView = models.ViewNetwork
			return m, nil
		case "6":
			m.currentView = models.ViewProcesses
			return m, nil
		case "7":
			m.currentView = models.ViewLogs
			return m, nil
		}

		switch m.currentView {
		case models.ViewProcesses:
			return m.updateProcessView(msg)
		case models.ViewLogs:
			return m.updateLogView(msg)
		}

	case TickMsg:
		if err := m.app.CollectMetrics(); err != nil {
			m.err = err
		}
		m.lastUpdate = time.Time(msg)
		return m, tickCmd()

	case error:
		m.err = msg
		return m, nil
	}

	return m, nil
}

func (m Model) updateProcessView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.processView.IsDialogActive() {
		if err := m.processView.HandleDialogInput(msg.String()); err != nil {
			m.err = err
		}
		return m, nil
	}

	if m.processView.IsSearching() {
		switch msg.String() {
		case "enter":
			m.processView.StopSearch()
		case "esc":
			m.processView.StopSearch()
		case "backspace":
			m.processView.HandleSearchBackspace()
		default:
			if len(msg.String()) == 1 {
				m.processView.HandleSearchInput(rune(msg.String()[0]))
			}
		}
		return m, nil
	}

	switch msg.String() {
	case "up", "k":
		m.processView.MoveUp()
	case "down", "j":
		m.processView.MoveDown()
	case "pgup":
		m.processView.PageUp()
	case "pgdown":
		m.processView.PageDown()
	case "/":
		m.processView.StartSearch()
	case "c":
		m.processView.SetSortField("cpu")
	case "m":
		m.processView.SetSortField("memory")
	case "n":
		m.processView.SetSortField("name")
	case "t":
		m.processView.SetSortField("time")
	case "s":
		m.processView.ToggleSortOrder()
	case "delete", "d":
		m.processView.ShowKillDialog()
	case "f":
		m.processView.ShowForceKillDialog()
	case "z":
		m.processView.ShowStopDialog()
	case "r":
		m.processView.ShowContinueDialog()
	case "P":
		m.processView.ShowNiceDialog()
	}
	return m, nil
}

func (m Model) updateLogView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.logView.MoveUp()
	case "down", "j":
		m.logView.MoveDown()
	case "pgup":
		m.logView.PageUp()
	case "pgdown":
		m.logView.PageDown()
	case "a":
		m.logView.ToggleAutoScroll()
	case "c":
		m.logView.ClearFilters()
	case "e":
		m.logView.SetLevelFilter("err")
	case "w":
		m.logView.SetLevelFilter("warn")
	case "i":
		m.logView.SetLevelFilter("info")
	}
	return m, nil
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	if m.showHelp {
		return m.renderHelp()
	}

	snapshot := m.app.GetLastSnapshot()
	header := m.renderHeader(snapshot)
	footer := m.renderFooter()

	headerHeight := lipgloss.Height(header)
	footerHeight := lipgloss.Height(footer)
	contentHeight := m.height - headerHeight - footerHeight
	if contentHeight < 0 {
		contentHeight = 0
	}

	var content string
	if snapshot == nil {
		content = lipgloss.NewStyle().
			Height(contentHeight).
			Width(m.width).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Collecting system metrics...")
	} else {
		viewContent := m.renderView(snapshot, contentHeight)
		content = lipgloss.NewStyle().Height(contentHeight).MaxHeight(contentHeight).Render(viewContent)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		footer,
	)
}

func (m Model) renderHeader(snapshot *models.MetricsSnapshot) string {
	title := "ltop - Linux System Monitor"

	var status string
	if snapshot != nil {
		status = fmt.Sprintf("Last Update: %s | %s",
			utils.FormatTime(snapshot.Timestamp),
			snapshot.Overview.Hostname)

		if m.app.GetState().Paused {
			status += " [PAUSED]"
		}
	}

	tabs := m.renderTabs()

	headerContent := lipgloss.JoinHorizontal(
		lipgloss.Left,
		styles.Title().Render(title),
		lipgloss.NewStyle().Margin(0, 2).Render("|"),
		status,
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerContent,
		tabs,
	)
}

func (m Model) renderTabs() string {
	tabs := []struct {
		key   string
		label string
		view  models.ViewType
	}{
		{"1", "Overview", models.ViewOverview},
		{"2", "CPU", models.ViewCPU},
		{"3", "Memory", models.ViewMemory},
		{"4", "Storage", models.ViewStorage},
		{"5", "Network", models.ViewNetwork},
		{"6", "Processes", models.ViewProcesses},
		{"7", "Logs", models.ViewLogs},
	}

	var tabStrings []string
	for _, tab := range tabs {
		style := styles.NavigationTab()
		if tab.view == m.currentView {
			style = styles.NavigationTabActive()
		}
		tabStr := fmt.Sprintf("[%s] %s", tab.key, tab.label)
		tabStrings = append(tabStrings, style.Render(tabStr))
	}

	return strings.Join(tabStrings, " ")
}

func (m Model) renderView(snapshot *models.MetricsSnapshot, height int) string {
	switch m.currentView {
	case models.ViewOverview:
		return m.overviewView.Render(snapshot, m.width, height)
	case models.ViewCPU:
		return m.cpuView.Render(snapshot, m.width, height)
	case models.ViewMemory:
		return m.memoryView.Render(snapshot, m.width, height)
	case models.ViewStorage:
		return m.storageView.Render(snapshot, m.width, height)
	case models.ViewNetwork:
		return m.networkView.Render(snapshot, m.width, height)
	case models.ViewProcesses:
		return m.processView.Render(snapshot, m.width, height)
	case models.ViewLogs:
		return m.logView.Render(snapshot, m.width, height)
	default:
		return "Unknown view"
	}
}

func (m Model) renderFooter() string {
	helpText := "Press 'h' for help, 'q' to quit, 'p' to pause, or 1-7 to switch views"
	switch m.currentView {
	case models.ViewLogs:
		helpText = "Logs: a=auto-scroll, c=clear filters, e/w/i=filter by error/warn/info"
	case models.ViewProcesses:
		if m.processView.IsDialogActive() {
			helpText = "Dialog: Enter=confirm, Esc=cancel, ←→=navigate"
		} else if m.processView.IsSearching() {
			helpText = "Search mode: Type to filter, Enter to apply, Esc to cancel"
		} else {
			helpText = "Processes: /=search, d=kill, f=force kill, z=stop, r=resume, P=priority"
		}
	}
	if m.err != nil {
		helpText = styles.Error().Render(fmt.Sprintf("Error: %v", m.err))
	}
	return styles.StatusBar().Width(m.width).Render(helpText)
}

func (m Model) renderHelp() string {
	help := `
ltop - Linux System Monitor Help

Navigation:
  1-7          Switch between views (Overview, CPU, Memory, Storage, Network, Processes, Logs)
  h, ?         Toggle this help screen
  q, Ctrl+C    Quit the application
  p            Pause/Resume updates
  r            Toggle refresh rate (1s/5s)
  T            Toggle theme (dark/light)

Process View (View 6):
  ↑/↓, k/j     Move selection up/down
  Page Up/Down Navigate by pages
  /            Start search/filter
  c/m/n/t      Sort by CPU/Memory/Name/Time
  s            Toggle sort order (asc/desc)
  d            Kill selected process (SIGTERM)
  f            Force kill selected process (SIGKILL)
  z            Stop selected process (SIGSTOP)
  r            Resume selected process (SIGCONT)
  P            Change process priority (nice)

Log View (View 7):
  ↑/↓, k/j     Move selection up/down
  Page Up/Down Navigate by pages
  a            Toggle auto-scroll (follows new entries)
  c            Clear all filters
  e            Filter by error level
  w            Filter by warning level
  i            Filter by info level

Views:
  1. Overview  - System summary with key metrics
  2. CPU       - Detailed CPU usage, load average, and per-core stats
  3. Memory    - RAM and swap usage with detailed breakdowns
  4. Storage   - Filesystem usage and disk I/O statistics
  5. Network   - Network interface statistics and bandwidth
  6. Processes - Process list with CPU, memory, and details
  7. Logs      - System logs with filtering and real-time monitoring

Configuration:
  Config file: ~/.config/ltop/config.json
  
Tips:
  - All metrics update every second by default
  - Use pause (p) to freeze the display for detailed inspection
  - Process list is sorted by CPU usage by default
  - Colors indicate status: green (normal), yellow (warning), red (critical)

Press 'h' again to return to the monitoring view.
`

	return styles.Border().
		Width(m.width - 4).
		Height(m.height - 4).
		Render(help)
}
