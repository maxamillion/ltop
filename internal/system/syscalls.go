package system

import (
	"syscall"
	"time"
	"unsafe"
)

type SystemInfo struct {
	Uptime    time.Duration
	Loads     [3]float64
	TotalRAM  uint64
	FreeRAM   uint64
	SharedRAM uint64
	BufferRAM uint64
	TotalSwap uint64
	FreeSwap  uint64
	Procs     uint16
}

type DiskUsage struct {
	Total uint64
	Free  uint64
	Used  uint64
}

func GetSystemInfo() (*SystemInfo, error) {
	var info syscall.Sysinfo_t
	err := syscall.Sysinfo(&info)
	if err != nil {
		return nil, err
	}

	return &SystemInfo{
		Uptime:    time.Duration(info.Uptime) * time.Second,
		Loads:     [3]float64{float64(info.Loads[0]) / 65536.0, float64(info.Loads[1]) / 65536.0, float64(info.Loads[2]) / 65536.0},
		TotalRAM:  info.Totalram * uint64(info.Unit),
		FreeRAM:   info.Freeram * uint64(info.Unit),
		SharedRAM: info.Sharedram * uint64(info.Unit),
		BufferRAM: info.Bufferram * uint64(info.Unit),
		TotalSwap: info.Totalswap * uint64(info.Unit),
		FreeSwap:  info.Freeswap * uint64(info.Unit),
		Procs:     info.Procs,
	}, nil
}

func GetDiskUsage(path string) (*DiskUsage, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return nil, err
	}

	total := uint64(stat.Blocks) * uint64(stat.Bsize)
	free := uint64(stat.Bavail) * uint64(stat.Bsize)
	used := total - free

	return &DiskUsage{
		Total: total,
		Free:  free,
		Used:  used,
	}, nil
}

type CPUTimes struct {
	User   uint64
	Nice   uint64
	System uint64
	Idle   uint64
	IOWait uint64
	IRQ    uint64
	SoftIRQ uint64
	Steal  uint64
	Guest  uint64
	GuestNice uint64
}

func GetProcessorCount() int {
	return int(syscall.SYS_SYSINFO)
}

type Rusage struct {
	Utime    syscall.Timeval
	Stime    syscall.Timeval
	Maxrss   int64
	Ixrss    int64
	Idrss    int64
	Isrss    int64
	Minflt   int64
	Majflt   int64
	Nswap    int64
	Inblock  int64
	Oublock  int64
	Msgsnd   int64
	Msgrcv   int64
	Nsignals int64
	Nvcsw    int64
	Nivcsw   int64
}

func GetResourceUsage() (*Rusage, error) {
	var ru syscall.Rusage
	err := syscall.Getrusage(syscall.RUSAGE_SELF, &ru)
	if err != nil {
		return nil, err
	}

	return (*Rusage)(unsafe.Pointer(&ru)), nil
}

func GetCurrentPID() int {
	return syscall.Getpid()
}

func GetParentPID() int {
	return syscall.Getppid()
}

func GetUID() int {
	return syscall.Getuid()
}

func GetGID() int {
	return syscall.Getgid()
}