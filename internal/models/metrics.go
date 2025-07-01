package models

import (
	"time"
)

type CPUMetrics struct {
	Usage       float64           `json:"usage"`
	LoadAverage [3]float64        `json:"load_average"`
	Cores       []CPUCoreMetrics  `json:"cores"`
	Frequency   map[string]uint64 `json:"frequency"`
	Temperature float64           `json:"temperature"`
	Times       CPUTimes          `json:"times"`
	Timestamp   time.Time         `json:"timestamp"`
}

type CPUCoreMetrics struct {
	ID    int      `json:"id"`
	Usage float64  `json:"usage"`
	Times CPUTimes `json:"times"`
}

type CPUTimes struct {
	User      uint64 `json:"user"`
	Nice      uint64 `json:"nice"`
	System    uint64 `json:"system"`
	Idle      uint64 `json:"idle"`
	IOWait    uint64 `json:"iowait"`
	IRQ       uint64 `json:"irq"`
	SoftIRQ   uint64 `json:"softirq"`
	Steal     uint64 `json:"steal"`
	Guest     uint64 `json:"guest"`
	GuestNice uint64 `json:"guest_nice"`
}

type MemoryMetrics struct {
	Total       uint64            `json:"total"`
	Free        uint64            `json:"free"`
	Available   uint64            `json:"available"`
	Used        uint64            `json:"used"`
	UsedPercent float64           `json:"used_percent"`
	Cached      uint64            `json:"cached"`
	Buffers     uint64            `json:"buffers"`
	Shared      uint64            `json:"shared"`
	Swap        SwapMetrics       `json:"swap"`
	Details     map[string]uint64 `json:"details"`
	Timestamp   time.Time         `json:"timestamp"`
}

type SwapMetrics struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

type StorageMetrics struct {
	Filesystems []FilesystemMetrics `json:"filesystems"`
	Disks       []DiskMetrics       `json:"disks"`
	IOStats     []DiskIOMetrics     `json:"io_stats"`
	Timestamp   time.Time           `json:"timestamp"`
}

type FilesystemMetrics struct {
	Device      string  `json:"device"`
	Mountpoint  string  `json:"mountpoint"`
	FSType      string  `json:"fstype"`
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
	InodesTotal uint64  `json:"inodes_total"`
	InodesFree  uint64  `json:"inodes_free"`
	InodesUsed  uint64  `json:"inodes_used"`
}

type DiskMetrics struct {
	Device    string `json:"device"`
	Model     string `json:"model"`
	Size      uint64 `json:"size"`
	ReadOnly  bool   `json:"read_only"`
	Removable bool   `json:"removable"`
}

type DiskIOMetrics struct {
	Device           string  `json:"device"`
	ReadIOs          uint64  `json:"read_ios"`
	ReadMerged       uint64  `json:"read_merged"`
	ReadSectors      uint64  `json:"read_sectors"`
	ReadTicks        uint64  `json:"read_ticks"`
	WriteIOs         uint64  `json:"write_ios"`
	WriteMerged      uint64  `json:"write_merged"`
	WriteSectors     uint64  `json:"write_sectors"`
	WriteTicks       uint64  `json:"write_ticks"`
	InFlight         uint64  `json:"in_flight"`
	IOTicks          uint64  `json:"io_ticks"`
	TimeInQueue      uint64  `json:"time_in_queue"`
	ReadBytesPerSec  float64 `json:"read_bytes_per_sec"`
	WriteBytesPerSec float64 `json:"write_bytes_per_sec"`
	IOPSRead         float64 `json:"iops_read"`
	IOPSWrite        float64 `json:"iops_write"`
	IOWaitPercent    float64 `json:"iowait_percent"`
}

type NetworkMetrics struct {
	Interfaces []NetworkInterface `json:"interfaces"`
	Timestamp  time.Time          `json:"timestamp"`
}

type NetworkInterface struct {
	Name            string  `json:"name"`
	BytesRecv       uint64  `json:"bytes_recv"`
	BytesSent       uint64  `json:"bytes_sent"`
	PacketsRecv     uint64  `json:"packets_recv"`
	PacketsSent     uint64  `json:"packets_sent"`
	ErrorsRecv      uint64  `json:"errors_recv"`
	ErrorsSent      uint64  `json:"errors_sent"`
	DroppedRecv     uint64  `json:"dropped_recv"`
	DroppedSent     uint64  `json:"dropped_sent"`
	Speed           uint64  `json:"speed"`
	Duplex          string  `json:"duplex"`
	MTU             int     `json:"mtu"`
	State           string  `json:"state"`
	RecvBytesPerSec float64 `json:"recv_bytes_per_sec"`
	SentBytesPerSec float64 `json:"sent_bytes_per_sec"`
}

type ProcessMetrics struct {
	Processes []Process `json:"processes"`
	Count     int       `json:"count"`
	Running   int       `json:"running"`
	Sleeping  int       `json:"sleeping"`
	Stopped   int       `json:"stopped"`
	Zombie    int       `json:"zombie"`
	Timestamp time.Time `json:"timestamp"`
}

type Process struct {
	PID           int            `json:"pid"`
	PPID          int            `json:"ppid"`
	Name          string         `json:"name"`
	Command       string         `json:"command"`
	State         string         `json:"state"`
	CPUPercent    float64        `json:"cpu_percent"`
	CPUTime       time.Duration  `json:"cpu_time"`
	MemoryRSS     uint64         `json:"memory_rss"`
	MemoryVMS     uint64         `json:"memory_vms"`
	MemoryPercent float64        `json:"memory_percent"`
	Priority      int            `json:"priority"`
	Nice          int            `json:"nice"`
	NumThreads    int            `json:"num_threads"`
	NumFDs        int            `json:"num_fds"`
	CreateTime    time.Time      `json:"create_time"`
	User          string         `json:"user"`
	Group         string         `json:"group"`
	IOStats       ProcessIOStats `json:"io_stats"`
	Children      []int          `json:"children"`
}

type ProcessIOStats struct {
	ReadBytes  uint64 `json:"read_bytes"`
	WriteBytes uint64 `json:"write_bytes"`
	ReadCount  uint64 `json:"read_count"`
	WriteCount uint64 `json:"write_count"`
}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Service   string    `json:"service"`
	Message   string    `json:"message"`
	Source    string    `json:"source"`
}

type LogMetrics struct {
	Entries    []LogEntry `json:"entries"`
	ErrorCount int        `json:"error_count"`
	WarnCount  int        `json:"warn_count"`
	InfoCount  int        `json:"info_count"`
	Sources    []string   `json:"sources"`
	Timestamp  time.Time  `json:"timestamp"`
}

type SystemOverview struct {
	Hostname     string        `json:"hostname"`
	Uptime       time.Duration `json:"uptime"`
	BootTime     time.Time     `json:"boot_time"`
	OS           string        `json:"os"`
	Platform     string        `json:"platform"`
	Kernel       string        `json:"kernel"`
	Architecture string        `json:"architecture"`
	CPUModel     string        `json:"cpu_model"`
	CPUCores     int           `json:"cpu_cores"`
	TotalMemory  uint64        `json:"total_memory"`
	CurrentUser  string        `json:"current_user"`
	LoadAvg      [3]float64    `json:"load_avg"`
}

type MetricsSnapshot struct {
	Overview  SystemOverview `json:"overview"`
	CPU       CPUMetrics     `json:"cpu"`
	Memory    MemoryMetrics  `json:"memory"`
	Storage   StorageMetrics `json:"storage"`
	Network   NetworkMetrics `json:"network"`
	Processes ProcessMetrics `json:"processes"`
	Logs      LogMetrics     `json:"logs"`
	Timestamp time.Time      `json:"timestamp"`
}
