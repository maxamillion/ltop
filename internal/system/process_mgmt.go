package system

import (
	"fmt"
	"os"
	"syscall"
)

type ProcessManager struct{}

func NewProcessManager() *ProcessManager {
	return &ProcessManager{}
}

func (pm *ProcessManager) KillProcess(pid int, signal syscall.Signal) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process %d: %w", pid, err)
	}

	err = process.Signal(signal)
	if err != nil {
		return fmt.Errorf("failed to send signal %v to process %d: %w", signal, pid, err)
	}

	return nil
}

func (pm *ProcessManager) TerminateProcess(pid int) error {
	return pm.KillProcess(pid, syscall.SIGTERM)
}

func (pm *ProcessManager) ForceKillProcess(pid int) error {
	return pm.KillProcess(pid, syscall.SIGKILL)
}

func (pm *ProcessManager) StopProcess(pid int) error {
	return pm.KillProcess(pid, syscall.SIGSTOP)
}

func (pm *ProcessManager) ContinueProcess(pid int) error {
	return pm.KillProcess(pid, syscall.SIGCONT)
}

func (pm *ProcessManager) SetProcessPriority(pid int, priority int) error {
	err := syscall.Setpriority(syscall.PRIO_PROCESS, pid, priority)
	if err != nil {
		return fmt.Errorf("failed to set priority for process %d: %w", pid, err)
	}
	return nil
}

func (pm *ProcessManager) GetProcessPriority(pid int) (int, error) {
	priority, err := syscall.Getpriority(syscall.PRIO_PROCESS, pid)
	if err != nil {
		return 0, fmt.Errorf("failed to get priority for process %d: %w", pid, err)
	}
	return priority, nil
}

func (pm *ProcessManager) CanManageProcess(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))
	return err == nil
}

func (pm *ProcessManager) IsProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))
	return err == nil
}
