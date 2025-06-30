package models

import (
	"time"
)

type SystemConfig struct {
	RefreshInterval time.Duration `json:"refresh_interval"`
	MaxProcesses    int           `json:"max_processes"`
	MaxLogEntries   int           `json:"max_log_entries"`
	LogSources      []string      `json:"log_sources"`
	Theme           string        `json:"theme"`
	ShowCPUPercent  bool          `json:"show_cpu_percent"`
	ShowMemoryPercent bool        `json:"show_memory_percent"`
	SortBy          string        `json:"sort_by"`
	SortOrder       string        `json:"sort_order"`
	ViewMode        string        `json:"view_mode"`
}

func DefaultSystemConfig() SystemConfig {
	return SystemConfig{
		RefreshInterval:   time.Second,
		MaxProcesses:      100,
		MaxLogEntries:     50,
		LogSources:        []string{"/var/log/syslog", "journalctl"},
		Theme:             "dark",
		ShowCPUPercent:    true,
		ShowMemoryPercent: true,
		SortBy:            "cpu",
		SortOrder:         "desc",
		ViewMode:          "overview",
	}
}

type AppState struct {
	CurrentView    string        `json:"current_view"`
	SelectedPID    int           `json:"selected_pid"`
	SearchQuery    string        `json:"search_query"`
	FilterActive   bool          `json:"filter_active"`
	Paused         bool          `json:"paused"`
	LastUpdate     time.Time     `json:"last_update"`
	ViewHistory    []string      `json:"view_history"`
	Errors         []AppError    `json:"errors"`
}

type AppError struct {
	Timestamp time.Time `json:"timestamp"`
	Component string    `json:"component"`
	Message   string    `json:"message"`
	Level     string    `json:"level"`
}

func (e AppError) Error() string {
	return e.Message
}

type ViewType string

const (
	ViewOverview  ViewType = "overview"
	ViewCPU       ViewType = "cpu"
	ViewMemory    ViewType = "memory"
	ViewStorage   ViewType = "storage"
	ViewNetwork   ViewType = "network"
	ViewProcesses ViewType = "processes"
	ViewLogs      ViewType = "logs"
)

type SortField string

const (
	SortByCPU    SortField = "cpu"
	SortByMemory SortField = "memory"
	SortByPID    SortField = "pid"
	SortByName   SortField = "name"
	SortByTime   SortField = "time"
)

type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

type ProcessState string

const (
	ProcessRunning   ProcessState = "R"
	ProcessSleeping  ProcessState = "S"
	ProcessStopped   ProcessState = "T"
	ProcessZombie    ProcessState = "Z"
	ProcessDiskSleep ProcessState = "D"
	ProcessIdle      ProcessState = "I"
)

func (ps ProcessState) String() string {
	switch ps {
	case ProcessRunning:
		return "Running"
	case ProcessSleeping:
		return "Sleeping"
	case ProcessStopped:
		return "Stopped"
	case ProcessZombie:
		return "Zombie"
	case ProcessDiskSleep:
		return "Disk Sleep"
	case ProcessIdle:
		return "Idle"
	default:
		return "Unknown"
	}
}

type LogLevel string

const (
	LogLevelEmerg  LogLevel = "emerg"
	LogLevelAlert  LogLevel = "alert"
	LogLevelCrit   LogLevel = "crit"
	LogLevelErr    LogLevel = "err"
	LogLevelWarn   LogLevel = "warn"
	LogLevelNotice LogLevel = "notice"
	LogLevelInfo   LogLevel = "info"
	LogLevelDebug  LogLevel = "debug"
)

func (ll LogLevel) Priority() int {
	switch ll {
	case LogLevelEmerg:
		return 0
	case LogLevelAlert:
		return 1
	case LogLevelCrit:
		return 2
	case LogLevelErr:
		return 3
	case LogLevelWarn:
		return 4
	case LogLevelNotice:
		return 5
	case LogLevelInfo:
		return 6
	case LogLevelDebug:
		return 7
	default:
		return 8
	}
}

func (ll LogLevel) String() string {
	return string(ll)
}

type ResourceThreshold struct {
	CPUWarning    float64 `json:"cpu_warning"`
	CPUCritical   float64 `json:"cpu_critical"`
	MemoryWarning float64 `json:"memory_warning"`
	MemoryCritical float64 `json:"memory_critical"`
	DiskWarning   float64 `json:"disk_warning"`
	DiskCritical  float64 `json:"disk_critical"`
	LoadWarning   float64 `json:"load_warning"`
	LoadCritical  float64 `json:"load_critical"`
}

func DefaultResourceThreshold() ResourceThreshold {
	return ResourceThreshold{
		CPUWarning:     70.0,
		CPUCritical:    90.0,
		MemoryWarning:  80.0,
		MemoryCritical: 95.0,
		DiskWarning:    85.0,
		DiskCritical:   95.0,
		LoadWarning:    2.0,
		LoadCritical:   5.0,
	}
}

type ColorTheme struct {
	Name       string `json:"name"`
	Background string `json:"background"`
	Foreground string `json:"foreground"`
	Primary    string `json:"primary"`
	Secondary  string `json:"secondary"`
	Success    string `json:"success"`
	Warning    string `json:"warning"`
	Error      string `json:"error"`
	Info       string `json:"info"`
	Muted      string `json:"muted"`
}

func DarkTheme() ColorTheme {
	return ColorTheme{
		Name:       "dark",
		Background: "#1a1a1a",
		Foreground: "#ffffff",
		Primary:    "#00ff00",
		Secondary:  "#ffff00",
		Success:    "#00ff00",
		Warning:    "#ffaa00",
		Error:      "#ff0000",
		Info:       "#00aaff",
		Muted:      "#666666",
	}
}

func LightTheme() ColorTheme {
	return ColorTheme{
		Name:       "light",
		Background: "#ffffff",
		Foreground: "#000000",
		Primary:    "#0066cc",
		Secondary:  "#666666",
		Success:    "#00cc00",
		Warning:    "#ff9900",
		Error:      "#cc0000",
		Info:       "#0099cc",
		Muted:      "#999999",
	}
}