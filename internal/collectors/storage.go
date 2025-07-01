package collectors

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/admiller/ltop/internal/models"
	"github.com/admiller/ltop/internal/system"
)

type StorageCollector struct {
	procReader    *system.ProcReader
	sysReader     *system.SysReader
	lastDiskStats map[string]models.DiskIOMetrics
	lastUpdate    time.Time
}

func NewStorageCollector() *StorageCollector {
	return &StorageCollector{
		procReader:    system.NewProcReader(),
		sysReader:     system.NewSysReader(),
		lastDiskStats: make(map[string]models.DiskIOMetrics),
		lastUpdate:    time.Now(),
	}
}

func (s *StorageCollector) Collect() (*models.StorageMetrics, error) {
	metrics := &models.StorageMetrics{
		Timestamp: time.Now(),
	}

	if err := s.collectFilesystems(metrics); err != nil {
		return nil, fmt.Errorf("failed to collect filesystems: %w", err)
	}

	if err := s.collectDiskStats(metrics); err != nil {
		return nil, fmt.Errorf("failed to collect disk stats: %w", err)
	}

	if err := s.collectBlockDevices(metrics); err != nil {
		return metrics, nil
	}

	return metrics, nil
}

func (s *StorageCollector) collectFilesystems(metrics *models.StorageMetrics) error {
	mounts, err := s.procReader.ReadMounts()
	if err != nil {
		return err
	}

	filesystems := make([]models.FilesystemMetrics, 0)
	seenMountpoints := make(map[string]bool)

	for _, mount := range mounts {
		fields := strings.Fields(mount)
		if len(fields) < 3 {
			continue
		}

		device := fields[0]
		mountpoint := fields[1]
		fstype := fields[2]

		if seenMountpoints[mountpoint] {
			continue
		}

		if !s.shouldIncludeFilesystem(device, mountpoint, fstype) {
			continue
		}

		diskUsage, err := system.GetDiskUsage(mountpoint)
		if err != nil {
			continue
		}

		fs := models.FilesystemMetrics{
			Device:      device,
			Mountpoint:  mountpoint,
			FSType:      fstype,
			Total:       diskUsage.Total,
			Free:        diskUsage.Free,
			Used:        diskUsage.Used,
			UsedPercent: float64(diskUsage.Used) / float64(diskUsage.Total) * 100.0,
		}

		filesystems = append(filesystems, fs)
		seenMountpoints[mountpoint] = true
	}

	metrics.Filesystems = filesystems
	return nil
}

func (s *StorageCollector) shouldIncludeFilesystem(device, mountpoint, fstype string) bool {
	if strings.HasPrefix(device, "/dev/loop") {
		return false
	}

	if strings.HasPrefix(mountpoint, "/proc") ||
		strings.HasPrefix(mountpoint, "/sys") ||
		strings.HasPrefix(mountpoint, "/dev") ||
		strings.HasPrefix(mountpoint, "/run") {
		return false
	}

	excludedFSTypes := []string{"tmpfs", "devtmpfs", "sysfs", "proc", "devpts", "cgroup", "cgroup2", "pstore", "bpf", "tracefs"}
	for _, excluded := range excludedFSTypes {
		if fstype == excluded {
			return false
		}
	}

	return true
}

func (s *StorageCollector) collectDiskStats(metrics *models.StorageMetrics) error {
	diskStats, err := s.procReader.ReadDiskStats()
	if err != nil {
		return err
	}

	ioStats := make([]models.DiskIOMetrics, 0)
	currentTime := time.Now()
	timeDelta := currentTime.Sub(s.lastUpdate).Seconds()

	for _, line := range diskStats {
		fields := strings.Fields(line)
		if len(fields) < 14 {
			continue
		}

		device := fields[2]

		if !s.shouldIncludeDisk(device) {
			continue
		}

		stat, err := s.parseDiskStatLine(fields)
		if err != nil {
			continue
		}

		stat.Device = device

		if lastStat, exists := s.lastDiskStats[device]; exists && timeDelta > 0 {
			s.calculateDiskRates(&stat, &lastStat, timeDelta)
		}

		ioStats = append(ioStats, stat)
		s.lastDiskStats[device] = stat
	}

	metrics.IOStats = ioStats
	s.lastUpdate = currentTime
	return nil
}

func (s *StorageCollector) shouldIncludeDisk(device string) bool {
	if len(device) < 2 {
		return false
	}

	if strings.HasPrefix(device, "loop") ||
		strings.HasPrefix(device, "ram") ||
		strings.HasPrefix(device, "dm-") {
		return false
	}

	lastChar := device[len(device)-1]
	if lastChar >= '0' && lastChar <= '9' {
		return false
	}

	return true
}

func (s *StorageCollector) parseDiskStatLine(fields []string) (models.DiskIOMetrics, error) {
	var stat models.DiskIOMetrics
	var err error

	if stat.ReadIOs, err = strconv.ParseUint(fields[3], 10, 64); err != nil {
		return stat, err
	}
	if stat.ReadMerged, err = strconv.ParseUint(fields[4], 10, 64); err != nil {
		return stat, err
	}
	if stat.ReadSectors, err = strconv.ParseUint(fields[5], 10, 64); err != nil {
		return stat, err
	}
	if stat.ReadTicks, err = strconv.ParseUint(fields[6], 10, 64); err != nil {
		return stat, err
	}
	if stat.WriteIOs, err = strconv.ParseUint(fields[7], 10, 64); err != nil {
		return stat, err
	}
	if stat.WriteMerged, err = strconv.ParseUint(fields[8], 10, 64); err != nil {
		return stat, err
	}
	if stat.WriteSectors, err = strconv.ParseUint(fields[9], 10, 64); err != nil {
		return stat, err
	}
	if stat.WriteTicks, err = strconv.ParseUint(fields[10], 10, 64); err != nil {
		return stat, err
	}
	if stat.InFlight, err = strconv.ParseUint(fields[11], 10, 64); err != nil {
		return stat, err
	}
	if stat.IOTicks, err = strconv.ParseUint(fields[12], 10, 64); err != nil {
		return stat, err
	}
	if stat.TimeInQueue, err = strconv.ParseUint(fields[13], 10, 64); err != nil {
		return stat, err
	}

	return stat, nil
}

func (s *StorageCollector) calculateDiskRates(current, last *models.DiskIOMetrics, timeDelta float64) {
	const sectorSize = 512

	readSectorsDelta := current.ReadSectors - last.ReadSectors
	writeSectorsDelta := current.WriteSectors - last.WriteSectors
	readIOsDelta := current.ReadIOs - last.ReadIOs
	writeIOsDelta := current.WriteIOs - last.WriteIOs

	current.ReadBytesPerSec = float64(readSectorsDelta*sectorSize) / timeDelta
	current.WriteBytesPerSec = float64(writeSectorsDelta*sectorSize) / timeDelta
	current.IOPSRead = float64(readIOsDelta) / timeDelta
	current.IOPSWrite = float64(writeIOsDelta) / timeDelta

	totalTicks := current.IOTicks - last.IOTicks
	current.IOWaitPercent = float64(totalTicks) / (timeDelta * 1000) * 100.0
	if current.IOWaitPercent > 100.0 {
		current.IOWaitPercent = 100.0
	}
}

func (s *StorageCollector) collectBlockDevices(metrics *models.StorageMetrics) error {
	devices, err := s.sysReader.ReadBlockDevices()
	if err != nil {
		return err
	}

	disks := make([]models.DiskMetrics, 0)

	for _, device := range devices {
		if !s.shouldIncludeDisk(device) {
			continue
		}

		disk := models.DiskMetrics{
			Device: device,
		}

		if size, err := s.sysReader.ReadBlockDeviceSize(device); err == nil {
			if sizeBlocks, err := strconv.ParseUint(size, 10, 64); err == nil {
				disk.Size = sizeBlocks * 512
			}
		}

		disks = append(disks, disk)
	}

	metrics.Disks = disks
	return nil
}
