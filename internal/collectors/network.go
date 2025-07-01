package collectors

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/admiller/ltop/internal/models"
	"github.com/admiller/ltop/internal/system"
)

type NetworkCollector struct {
	procReader         *system.ProcReader
	sysReader          *system.SysReader
	lastInterfaceStats map[string]models.NetworkInterface
	lastUpdate         time.Time
}

func NewNetworkCollector() *NetworkCollector {
	return &NetworkCollector{
		procReader:         system.NewProcReader(),
		sysReader:          system.NewSysReader(),
		lastInterfaceStats: make(map[string]models.NetworkInterface),
		lastUpdate:         time.Now(),
	}
}

func (n *NetworkCollector) Collect() (*models.NetworkMetrics, error) {
	metrics := &models.NetworkMetrics{
		Timestamp: time.Now(),
	}

	if err := n.collectNetworkStats(metrics); err != nil {
		return nil, fmt.Errorf("failed to collect network stats: %w", err)
	}

	return metrics, nil
}

func (n *NetworkCollector) collectNetworkStats(metrics *models.NetworkMetrics) error {
	netStats, err := n.procReader.ReadNetworkStats()
	if err != nil {
		return err
	}

	interfaces := make([]models.NetworkInterface, 0)
	currentTime := time.Now()
	timeDelta := currentTime.Sub(n.lastUpdate).Seconds()

	for i, line := range netStats {
		if i < 2 {
			continue
		}

		iface, err := n.parseNetworkLine(line)
		if err != nil {
			continue
		}

		if !n.shouldIncludeInterface(iface.Name) {
			continue
		}

		if err := n.collectInterfaceDetails(&iface); err != nil {
			continue
		}

		if lastIface, exists := n.lastInterfaceStats[iface.Name]; exists && timeDelta > 0 {
			n.calculateNetworkRates(&iface, &lastIface, timeDelta)
		}

		interfaces = append(interfaces, iface)
		n.lastInterfaceStats[iface.Name] = iface
	}

	metrics.Interfaces = interfaces
	n.lastUpdate = currentTime
	return nil
}

func (n *NetworkCollector) parseNetworkLine(line string) (models.NetworkInterface, error) {
	var iface models.NetworkInterface

	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return iface, fmt.Errorf("invalid network line format")
	}

	iface.Name = strings.TrimSpace(parts[0])

	fields := strings.Fields(parts[1])
	if len(fields) < 16 {
		return iface, fmt.Errorf("insufficient network fields")
	}

	var err error
	if iface.BytesRecv, err = strconv.ParseUint(fields[0], 10, 64); err != nil {
		return iface, err
	}
	if iface.PacketsRecv, err = strconv.ParseUint(fields[1], 10, 64); err != nil {
		return iface, err
	}
	if iface.ErrorsRecv, err = strconv.ParseUint(fields[2], 10, 64); err != nil {
		return iface, err
	}
	if iface.DroppedRecv, err = strconv.ParseUint(fields[3], 10, 64); err != nil {
		return iface, err
	}

	if iface.BytesSent, err = strconv.ParseUint(fields[8], 10, 64); err != nil {
		return iface, err
	}
	if iface.PacketsSent, err = strconv.ParseUint(fields[9], 10, 64); err != nil {
		return iface, err
	}
	if iface.ErrorsSent, err = strconv.ParseUint(fields[10], 10, 64); err != nil {
		return iface, err
	}
	if iface.DroppedSent, err = strconv.ParseUint(fields[11], 10, 64); err != nil {
		return iface, err
	}

	return iface, nil
}

func (n *NetworkCollector) shouldIncludeInterface(name string) bool {
	if name == "lo" {
		return false
	}

	excludePrefixes := []string{"docker", "br-", "veth", "virbr", "tap"}
	for _, prefix := range excludePrefixes {
		if strings.HasPrefix(name, prefix) {
			return false
		}
	}

	return true
}

func (n *NetworkCollector) collectInterfaceDetails(iface *models.NetworkInterface) error {
	if state, err := n.sysReader.ReadNetworkOperState(iface.Name); err == nil {
		iface.State = state
	}

	if speed, err := n.sysReader.ReadNetworkSpeed(iface.Name); err == nil {
		if speedMbps, err := strconv.ParseUint(speed, 10, 64); err == nil {
			iface.Speed = speedMbps * 1000000
		}
	}

	return nil
}

func (n *NetworkCollector) calculateNetworkRates(current, last *models.NetworkInterface, timeDelta float64) {
	bytesRecvDelta := current.BytesRecv - last.BytesRecv
	bytesSentDelta := current.BytesSent - last.BytesSent

	current.RecvBytesPerSec = float64(bytesRecvDelta) / timeDelta
	current.SentBytesPerSec = float64(bytesSentDelta) / timeDelta
}
