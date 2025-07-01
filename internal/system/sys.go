package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SysReader struct {
	basePath string
}

func NewSysReader() *SysReader {
	return &SysReader{
		basePath: "/sys",
	}
}

func (s *SysReader) ReadFile(path string) ([]byte, error) {
	fullPath := filepath.Join(s.basePath, path)
	return os.ReadFile(fullPath)
}

func (s *SysReader) ReadString(path string) (string, error) {
	data, err := s.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func (s *SysReader) FileExists(path string) bool {
	fullPath := filepath.Join(s.basePath, path)
	_, err := os.Stat(fullPath)
	return err == nil
}

func (s *SysReader) ListDir(path string) ([]string, error) {
	fullPath := filepath.Join(s.basePath, path)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names, nil
}

func (s *SysReader) ReadCPUFreq(cpu string) (string, error) {
	path := fmt.Sprintf("devices/system/cpu/cpu%s/cpufreq/scaling_cur_freq", cpu)
	if !s.FileExists(path) {
		return "", nil // Return empty string if file doesn't exist, no error
	}
	return s.ReadString(path)
}

func (s *SysReader) ReadCPUMaxFreq(cpu string) (string, error) {
	path := fmt.Sprintf("devices/system/cpu/cpu%s/cpufreq/scaling_max_freq", cpu)
	return s.ReadString(path)
}

func (s *SysReader) ReadCPUTemperature() (string, error) {
	thermalZones, err := s.ListDir("class/thermal")
	if err != nil {
		return "", err
	}

	for _, zone := range thermalZones {
		if strings.HasPrefix(zone, "thermal_zone") {
			tempPath := fmt.Sprintf("class/thermal/%s/temp", zone)
			if s.FileExists(tempPath) {
				return s.ReadString(tempPath)
			}
		}
	}
	return "", fmt.Errorf("no thermal zones found")
}

func (s *SysReader) ReadBlockDevices() ([]string, error) {
	return s.ListDir("block")
}

func (s *SysReader) ReadBlockDeviceStat(device string) (string, error) {
	path := fmt.Sprintf("block/%s/stat", device)
	return s.ReadString(path)
}

func (s *SysReader) ReadBlockDeviceSize(device string) (string, error) {
	path := fmt.Sprintf("block/%s/size", device)
	return s.ReadString(path)
}

func (s *SysReader) ReadNetworkInterfaces() ([]string, error) {
	return s.ListDir("class/net")
}

func (s *SysReader) ReadNetworkSpeed(iface string) (string, error) {
	path := fmt.Sprintf("class/net/%s/speed", iface)
	if s.FileExists(path) {
		return s.ReadString(path)
	}
	return "", fmt.Errorf("speed not available for interface %s", iface)
}

func (s *SysReader) ReadNetworkOperState(iface string) (string, error) {
	path := fmt.Sprintf("class/net/%s/operstate", iface)
	return s.ReadString(path)
}

func (s *SysReader) ReadMemoryInfo() (map[string]string, error) {
	result := make(map[string]string)

	memInfo := []string{
		"kernel/mm/transparent_hugepage/enabled",
		"kernel/mm/ksm/run",
	}

	for _, path := range memInfo {
		if s.FileExists(path) {
			value, err := s.ReadString(path)
			if err == nil {
				result[path] = value
			}
		}
	}

	return result, nil
}

func (s *SysReader) ReadPowerSupplyInfo() (map[string]string, error) {
	result := make(map[string]string)

	supplies, err := s.ListDir("class/power_supply")
	if err != nil {
		return result, nil
	}

	for _, supply := range supplies {
		for _, prop := range []string{"capacity", "status", "present"} {
			path := fmt.Sprintf("class/power_supply/%s/%s", supply, prop)
			if s.FileExists(path) {
				value, err := s.ReadString(path)
				if err == nil {
					key := fmt.Sprintf("%s_%s", supply, prop)
					result[key] = value
				}
			}
		}
	}

	return result, nil
}
