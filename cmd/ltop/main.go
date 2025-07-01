package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/admiller/ltop/internal/app"
	"github.com/admiller/ltop/internal/models"
	"github.com/admiller/ltop/internal/ui/views"
	"github.com/admiller/ltop/pkg/utils"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	AppName    = "ltop"
	AppVersion = "0.1.0"
	AppDesc    = "Linux system monitoring tool with modern TUI"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-v", "--version":
			fmt.Printf("%s version %s\n", AppName, AppVersion)
			fmt.Printf("%s\n", AppDesc)
			return
		case "-h", "--help":
			printHelp()
			return
		}
	}

	ltopApp := app.New()

	if err := ltopApp.LoadConfig(); err != nil {
		log.Printf("Failed to load config, using defaults: %v", err)
	}

	if len(os.Args) > 1 && os.Args[1] == "--demo" {
		runDemo(ltopApp)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "--tui" {
		log.Printf("Starting %s %s in TUI mode", AppName, AppVersion)
		if err := runTUI(ltopApp); err != nil {
			log.Fatalf("TUI application failed: %v", err)
		}
		return
	}

	log.Printf("Starting %s %s in TUI mode", AppName, AppVersion)
	if err := runTUI(ltopApp); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}

func printHelp() {
	fmt.Printf("%s - %s\n\n", AppName, AppDesc)
	fmt.Println("Usage:")
	fmt.Printf("  %s [options]\n\n", AppName)
	fmt.Println("Options:")
	fmt.Println("  -h, --help     Show this help message")
	fmt.Println("  -v, --version  Show version information")
	fmt.Println("  --demo         Run demo mode (command line output)")
	fmt.Println("  --tui          Run TUI mode (default)")
	fmt.Println("")
	fmt.Println("Interactive Commands:")
	fmt.Println("  q, Ctrl+C      Quit")
	fmt.Println("  p              Pause/Resume updates")
	fmt.Println("  1-7            Switch between views")
	fmt.Println("  ↑↓, k/j        Navigate lists")
	fmt.Println("  a              Toggle auto-scroll (logs)")
	fmt.Println("  e/w/i          Filter logs by level")
	fmt.Println("")
	fmt.Println("Views:")
	fmt.Println("  1 - System Overview")
	fmt.Println("  2 - CPU Monitoring")
	fmt.Println("  3 - Memory Usage")
	fmt.Println("  4 - Storage & I/O")
	fmt.Println("  5 - Network Interfaces")
	fmt.Println("  6 - Process List")
	fmt.Println("  7 - System Logs")
	fmt.Println("")
	fmt.Printf("For more information, visit: https://github.com/admiller/ltop\n")
}

func runDemo(ltopApp *app.App) {
	fmt.Printf("=== %s Demo Mode ===\n\n", AppName)

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for i := 0; i < 10; i++ {
			if err := ltopApp.CollectMetrics(); err != nil {
				log.Printf("Metrics collection failed: %v", err)
			}
			<-ticker.C
		}
	}()

	time.Sleep(2 * time.Second)

	for i := 0; i < 3; i++ {
		snapshot := ltopApp.GetLastSnapshot()
		if snapshot == nil {
			fmt.Println("Waiting for metrics...")
			time.Sleep(time.Second)
			continue
		}

		printSystemInfo(snapshot)
		if i < 2 {
			time.Sleep(3 * time.Second)
			fmt.Println(strings.Repeat("-", 80))
		}
	}
}

func printSystemInfo(snapshot *models.MetricsSnapshot) {
	fmt.Printf("Timestamp: %s\n", utils.FormatDateTime(snapshot.Timestamp))
	fmt.Printf("Hostname: %s | User: %s\n", snapshot.Overview.Hostname, snapshot.Overview.CurrentUser)
	fmt.Printf("Uptime: %s\n", utils.FormatUptime(snapshot.Overview.Uptime))
	fmt.Println()

	fmt.Printf("CPU Usage: %s | Load: %s\n",
		utils.FormatPercent(snapshot.CPU.Usage),
		utils.FormatLoadAverage(snapshot.CPU.LoadAverage))

	if snapshot.CPU.Temperature != 0 {
		fmt.Printf("CPU Temperature: %s\n", utils.FormatTemperature(snapshot.CPU.Temperature))
	}
	fmt.Println()

	fmt.Printf("Memory: %s / %s (%s)\n",
		utils.FormatBytes(snapshot.Memory.Used),
		utils.FormatBytes(snapshot.Memory.Total),
		utils.FormatPercent(snapshot.Memory.UsedPercent))

	if snapshot.Memory.Swap.Total > 0 {
		fmt.Printf("Swap: %s / %s (%s)\n",
			utils.FormatBytes(snapshot.Memory.Swap.Used),
			utils.FormatBytes(snapshot.Memory.Swap.Total),
			utils.FormatPercent(snapshot.Memory.Swap.UsedPercent))
	}
	fmt.Println()

	fmt.Printf("Processes: %d (Running: %d, Sleeping: %d)\n",
		snapshot.Processes.Count,
		snapshot.Processes.Running,
		snapshot.Processes.Sleeping)

	if len(snapshot.Storage.Filesystems) > 0 {
		fmt.Println("\nStorage:")
		for _, fs := range snapshot.Storage.Filesystems {
			if len(fs.Mountpoint) > 1 {
				fmt.Printf("  %s: %s / %s (%s)\n",
					utils.TruncateString(fs.Mountpoint, 20),
					utils.FormatBytes(fs.Used),
					utils.FormatBytes(fs.Total),
					utils.FormatPercent(fs.UsedPercent))
			}
		}
	}

	if len(snapshot.Network.Interfaces) > 0 {
		fmt.Println("\nNetwork Interfaces:")
		for _, iface := range snapshot.Network.Interfaces {
			fmt.Printf("  %s: ↓%s ↑%s (%s)\n",
				iface.Name,
				utils.FormatBytesPerSecond(iface.RecvBytesPerSec),
				utils.FormatBytesPerSecond(iface.SentBytesPerSec),
				iface.State)
		}
	}

	if len(snapshot.Processes.Processes) > 0 {
		fmt.Println("\nTop 5 Processes by CPU:")
		fmt.Printf("%-8s %-20s %-10s %-10s %s\n", "PID", "NAME", "CPU%", "MEMORY", "COMMAND")

		count := 5
		if len(snapshot.Processes.Processes) < count {
			count = len(snapshot.Processes.Processes)
		}

		for i := 0; i < count; i++ {
			proc := snapshot.Processes.Processes[i]
			fmt.Printf("%-8d %-20s %-10s %-10s %s\n",
				proc.PID,
				utils.TruncateString(proc.Name, 20),
				utils.FormatPercent(proc.CPUPercent),
				utils.FormatBytes(proc.MemoryRSS),
				utils.TruncateString(proc.Command, 40))
		}
	}
}

func runTUI(ltopApp *app.App) error {
	model := views.NewModel(ltopApp)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	if err != nil {
		log.Printf("TUI error: %v", err)
		return err
	}

	return nil
}
