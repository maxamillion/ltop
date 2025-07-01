package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/admiller/ltop/internal/collectors"
	"github.com/admiller/ltop/internal/models"
)

type App struct {
	config           models.SystemConfig
	state            models.AppState
	cpuCollector     *collectors.CPUCollector
	memoryCollector  *collectors.MemoryCollector
	processCollector *collectors.ProcessCollector
	storageCollector *collectors.StorageCollector
	networkCollector *collectors.NetworkCollector
	logCollector     *collectors.LogCollector
	ctx              context.Context
	cancel           context.CancelFunc
	lastSnapshot     *models.MetricsSnapshot
}

func New() *App {
	ctx, cancel := context.WithCancel(context.Background())
	config := models.DefaultSystemConfig()

	return &App{
		config:           config,
		state:            models.AppState{},
		cpuCollector:     collectors.NewCPUCollector(),
		memoryCollector:  collectors.NewMemoryCollector(),
		processCollector: collectors.NewProcessCollector(),
		storageCollector: collectors.NewStorageCollector(),
		networkCollector: collectors.NewNetworkCollector(),
		logCollector:     collectors.NewLogCollector(config.LogSources, config.MaxLogEntries),
		ctx:              ctx,
		cancel:           cancel,
	}
}

func (a *App) Run() error {
	defer a.cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(a.config.RefreshInterval)
	defer ticker.Stop()

	if err := a.CollectMetrics(); err != nil {
		log.Printf("Initial metrics collection failed: %v", err)
	}

	for {
		select {
		case <-a.ctx.Done():
			return nil
		case <-sigChan:
			log.Println("Received shutdown signal")
			return nil
		case <-ticker.C:
			if !a.state.Paused {
				if err := a.CollectMetrics(); err != nil {
					log.Printf("Metrics collection failed: %v", err)
				}
			}
		}
	}
}

func (a *App) CollectMetrics() error {
	start := time.Now()

	snapshot := &models.MetricsSnapshot{
		Timestamp: start,
	}

	if err := a.collectSystemOverview(snapshot); err != nil {
		return err
	}

	if cpuMetrics, err := a.cpuCollector.Collect(); err == nil {
		snapshot.CPU = *cpuMetrics
	} else {
		log.Printf("CPU collection failed: %v", err)
	}

	if memoryMetrics, err := a.memoryCollector.Collect(); err == nil {
		snapshot.Memory = *memoryMetrics
	} else {
		log.Printf("Memory collection failed: %v", err)
	}

	if processMetrics, err := a.processCollector.Collect(); err == nil {
		snapshot.Processes = *processMetrics
	} else {
		log.Printf("Process collection failed: %v", err)
	}

	if storageMetrics, err := a.storageCollector.Collect(); err == nil {
		snapshot.Storage = *storageMetrics
	} else {
		log.Printf("Storage collection failed: %v", err)
	}

	if networkMetrics, err := a.networkCollector.Collect(); err == nil {
		snapshot.Network = *networkMetrics
	} else {
		log.Printf("Network collection failed: %v", err)
	}

	if logMetrics, err := a.logCollector.Collect(); err == nil {
		snapshot.Logs = *logMetrics
	} else {
		log.Printf("Log collection failed: %v", err)
	}

	a.lastSnapshot = snapshot
	a.state.LastUpdate = time.Now()

	return nil
}

func (a *App) collectSystemOverview(snapshot *models.MetricsSnapshot) error {
	hostname, _ := os.Hostname()

	snapshot.Overview = models.SystemOverview{
		Hostname:    hostname,
		CurrentUser: os.Getenv("USER"),
	}

	return nil
}

func (a *App) GetLastSnapshot() *models.MetricsSnapshot {
	return a.lastSnapshot
}

func (a *App) GetConfig() models.SystemConfig {
	return a.config
}

func (a *App) SetConfig(config models.SystemConfig) {
	a.config = config
}

func (a *App) GetState() models.AppState {
	return a.state
}

func (a *App) SetCurrentView(view string) {
	a.state.CurrentView = view
}

func (a *App) TogglePause() {
	a.state.Paused = !a.state.Paused
}

func (a *App) SetSearchQuery(query string) {
	a.state.SearchQuery = query
}

func (a *App) SetSelectedPID(pid int) {
	a.state.SelectedPID = pid
}

func (a *App) Shutdown() {
	a.cancel()
}
