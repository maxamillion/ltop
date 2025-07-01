package system

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type ProcReader struct {
	basePath string
}

func NewProcReader() *ProcReader {
	return &ProcReader{
		basePath: "/proc",
	}
}

func (p *ProcReader) ReadFile(path string) ([]byte, error) {
	fullPath := fmt.Sprintf("%s/%s", p.basePath, path)
	return os.ReadFile(fullPath)
}

func (p *ProcReader) ReadLines(path string) ([]string, error) {
	fullPath := fmt.Sprintf("%s/%s", p.basePath, path)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func (p *ProcReader) ReadFirstLine(path string) (string, error) {
	fullPath := fmt.Sprintf("%s/%s", p.basePath, path)
	file, err := os.Open(fullPath)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", scanner.Err()
}

func (p *ProcReader) ReadKeyValuePairs(path string) (map[string]string, error) {
	lines, err := p.ReadLines(path)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = value
		}
	}
	return result, nil
}

func (p *ProcReader) ReadCPUStat() ([]string, error) {
	return p.ReadLines("stat")
}

func (p *ProcReader) ReadMemInfo() (map[string]string, error) {
	return p.ReadKeyValuePairs("meminfo")
}

func (p *ProcReader) ReadLoadAvg() (string, error) {
	return p.ReadFirstLine("loadavg")
}

func (p *ProcReader) ReadUptime() (string, error) {
	return p.ReadFirstLine("uptime")
}

func (p *ProcReader) ReadProcesses() ([]string, error) {
	entries, err := os.ReadDir(p.basePath)
	if err != nil {
		return nil, err
	}

	var pids []string
	for _, entry := range entries {
		if entry.IsDir() {
			if _, err := strconv.Atoi(entry.Name()); err == nil {
				pids = append(pids, entry.Name())
			}
		}
	}
	return pids, nil
}

func (p *ProcReader) ReadProcessStat(pid string) (string, error) {
	return p.ReadFirstLine(fmt.Sprintf("%s/stat", pid))
}

func (p *ProcReader) ReadProcessStatus(pid string) (map[string]string, error) {
	return p.ReadKeyValuePairs(fmt.Sprintf("%s/status", pid))
}

func (p *ProcReader) ReadProcessCmdline(pid string) (string, error) {
	data, err := p.ReadFile(fmt.Sprintf("%s/cmdline", pid))
	if err != nil {
		return "", err
	}

	cmdline := string(data)
	cmdline = strings.ReplaceAll(cmdline, "\x00", " ")
	return strings.TrimSpace(cmdline), nil
}

func (p *ProcReader) ReadDiskStats() ([]string, error) {
	return p.ReadLines("diskstats")
}

func (p *ProcReader) ReadNetworkStats() ([]string, error) {
	return p.ReadLines("net/dev")
}

func (p *ProcReader) ReadMounts() ([]string, error) {
	return p.ReadLines("mounts")
}

func (p *ProcReader) OpenFile(path string) (io.ReadCloser, error) {
	fullPath := fmt.Sprintf("%s/%s", p.basePath, path)
	return os.Open(fullPath)
}

func (p *ProcReader) FileExists(path string) bool {
	fullPath := fmt.Sprintf("%s/%s", p.basePath, path)
	_, err := os.Stat(fullPath)
	return err == nil
}
