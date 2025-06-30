# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Essential Commands

### Build and Test
```bash
make build       # Build the ltop binary to bin/ltop
make test        # Run all tests
make test-coverage # Run tests with HTML coverage report
make check       # Run fmt, vet, lint, and test in sequence
```

### Development
```bash
make run         # Run directly with go run
make dev         # Build and run the binary
go test ./internal/collectors -v  # Run tests for specific package
```

### Code Quality
```bash
make lint        # Run golangci-lint (requires installation)
make fmt         # Format all Go code
make vet         # Run go vet
make security    # Run gosec security checks
```

## Architecture Overview

ltop is a Linux system monitoring tool built with Go and Bubble Tea TUI framework. The architecture follows a layered approach:

### Core Components

**Data Collection Layer** (`internal/collectors/`)
- Each collector (CPU, Memory, Storage, Network, Processes, Logs) implements independent metric gathering
- Collectors read from Linux `/proc` and `/sys` filesystems via `internal/system/` interfaces
- All collectors return structured data via `internal/models/MetricsSnapshot`

**System Interface Layer** (`internal/system/`)
- `proc.go`: Handles `/proc` filesystem reads (CPU stats, memory, processes)
- `sys.go`: Handles `/sys` filesystem reads (hardware info, temperatures)
- `process_mgmt.go`: Process control operations (kill, stop, priority changes)
- `syscalls.go`: Low-level system call wrappers

**UI Layer** (`internal/ui/`)
- Built on Bubble Tea framework with 7 specialized views (Overview, CPU, Memory, Storage, Network, Processes, Logs)
- Views are stateful and handle their own navigation, sorting, and filtering
- `components/`: Reusable TUI components (tables, dialogs, inputs, gauges)
- `styles/`: Centralized theming and styling

**Application Layer** (`internal/app/`)
- Orchestrates data collection and UI coordination
- Manages application lifecycle and configuration
- Handles real-time metric updates and state management

### Key Design Patterns

**Collector Pattern**: Each system metric type has its own collector that implements a `Collect()` method returning typed metrics.

**View Pattern**: Each UI view is self-contained with its own state, rendering logic, and input handling.

**Configuration-Driven**: System behavior controlled via `models.SystemConfig` with sensible defaults.

### Data Flow
1. App coordinates periodic collection from all collectors
2. Collectors read system data via system interfaces
3. Raw data is structured into `MetricsSnapshot`
4. UI views render snapshot data with Bubble Tea
5. User interactions trigger view state changes or system operations

### Testing Strategy
- Unit tests for collectors validate system data parsing
- Utility function tests ensure formatting and math operations
- Focus on `/proc` and `/sys` data parsing accuracy
- Process management operations tested with mock data

### Process Management
The tool provides interactive process control through the Processes view:
- Kill/terminate processes with confirmation dialogs
- Stop/continue processes (SIGSTOP/SIGCONT)
- Adjust process priority (nice values)
- All operations require user confirmation for safety