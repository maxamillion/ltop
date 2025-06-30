package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/admiller/ltop/internal/models"
)

const (
	DefaultConfigDir  = ".config/ltop"
	DefaultConfigFile = "config.json"
)

func (a *App) LoadConfig() error {
	configPath, err := a.getConfigPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return a.SaveConfig()
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var config models.SystemConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	a.config = config
	return nil
}

func (a *App) SaveConfig() error {
	configPath, err := a.getConfigPath()
	if err != nil {
		return err
	}

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(a.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func (a *App) getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, DefaultConfigDir, DefaultConfigFile), nil
}

func (a *App) UpdateRefreshInterval(interval time.Duration) {
	a.config.RefreshInterval = interval
}

func (a *App) UpdateMaxProcesses(max int) {
	a.config.MaxProcesses = max
}

func (a *App) UpdateTheme(theme string) {
	a.config.Theme = theme
}

func (a *App) UpdateSortBy(sortBy string) {
	a.config.SortBy = sortBy
}

func (a *App) UpdateSortOrder(order string) {
	a.config.SortOrder = order
}

func (a *App) UpdateViewMode(mode string) {
	a.config.ViewMode = mode
}

func (a *App) ToggleCPUPercent() {
	a.config.ShowCPUPercent = !a.config.ShowCPUPercent
}

func (a *App) ToggleMemoryPercent() {
	a.config.ShowMemoryPercent = !a.config.ShowMemoryPercent
}