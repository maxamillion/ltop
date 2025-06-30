package views

import (
	"fmt"
	"strings"

	"github.com/admiller/ltop/internal/models"
	"github.com/admiller/ltop/internal/ui/components"
	"github.com/admiller/ltop/internal/ui/styles"
	"github.com/admiller/ltop/pkg/utils"
)

type StorageView struct {
	fsGauges []*components.Gauge
}

func NewStorageView() *StorageView {
	return &StorageView{
		fsGauges: make([]*components.Gauge, 0),
	}
}

func (sv *StorageView) Render(snapshot *models.MetricsSnapshot, width, height int) string {
	if snapshot == nil {
		return "No data available"
	}

	var sections []string

	sections = append(sections, sv.renderFilesystems(snapshot))
	sections = append(sections, sv.renderDiskIO(snapshot))

	content := strings.Join(sections, "\n\n")
	return styles.Panel().Width(width).Render(content)
}

func (sv *StorageView) renderFilesystems(snapshot *models.MetricsSnapshot) string {
	var fs []string
	fs = append(fs, styles.Title().Render("Filesystem Usage"))

	if len(snapshot.Storage.Filesystems) == 0 {
		fs = append(fs, styles.Muted().Render("No filesystems found"))
		return strings.Join(fs, "\n")
	}

	for len(sv.fsGauges) < len(snapshot.Storage.Filesystems) {
		sv.fsGauges = append(sv.fsGauges, components.NewGauge(30))
	}

	for i, filesystem := range snapshot.Storage.Filesystems {
		if i < len(sv.fsGauges) {
			mountpoint := utils.TruncateString(filesystem.Mountpoint, 20)
			gauge := sv.fsGauges[i].Render(filesystem.UsedPercent, mountpoint)
			fs = append(fs, gauge)

			detail := fmt.Sprintf("  %s / %s (%s) on %s",
				utils.FormatBytes(filesystem.Used),
				utils.FormatBytes(filesystem.Total),
				filesystem.FSType,
				filesystem.Device)
			fs = append(fs, styles.Muted().Render(detail))
		}
	}

	return strings.Join(fs, "\n")
}

func (sv *StorageView) renderDiskIO(snapshot *models.MetricsSnapshot) string {
	var io []string
	io = append(io, styles.Title().Render("Disk I/O Statistics"))

	if len(snapshot.Storage.IOStats) == 0 {
		io = append(io, styles.Muted().Render("No disk I/O statistics available"))
		return strings.Join(io, "\n")
	}

	headers := []string{"Device", "Read/s", "Write/s", "IOPS Read", "IOPS Write", "I/O Wait%"}
	io = append(io, sv.renderIOHeader(headers))

	for _, stat := range snapshot.Storage.IOStats {
		row := sv.renderIORow(stat)
		io = append(io, row)
	}

	return strings.Join(io, "\n")
}

func (sv *StorageView) renderIOHeader(headers []string) string {
	var parts []string
	widths := []int{12, 12, 12, 12, 12, 12}

	for i, header := range headers {
		if i < len(widths) {
			text := utils.PadString(header, widths[i], ' ')
			parts = append(parts, styles.TableHeader().Render(text))
		}
	}

	return strings.Join(parts, " ")
}

func (sv *StorageView) renderIORow(stat models.DiskIOMetrics) string {
	var parts []string
	widths := []int{12, 12, 12, 12, 12, 12}

	device := utils.TruncateString(stat.Device, widths[0])
	device = utils.PadString(device, widths[0], ' ')
	parts = append(parts, styles.TableRow().Render(device))

	readRate := utils.PadString(utils.FormatBytesPerSecond(stat.ReadBytesPerSec), widths[1], ' ')
	parts = append(parts, styles.TableRow().Render(readRate))

	writeRate := utils.PadString(utils.FormatBytesPerSecond(stat.WriteBytesPerSec), widths[2], ' ')
	parts = append(parts, styles.TableRow().Render(writeRate))

	iopsRead := utils.PadString(fmt.Sprintf("%.1f", stat.IOPSRead), widths[3], ' ')
	parts = append(parts, styles.TableRow().Render(iopsRead))

	iopsWrite := utils.PadString(fmt.Sprintf("%.1f", stat.IOPSWrite), widths[4], ' ')
	parts = append(parts, styles.TableRow().Render(iopsWrite))

	iowait := utils.PadString(utils.FormatPercent(stat.IOWaitPercent), widths[5], ' ')
	iowaitStyle := styles.PercentageColor(stat.IOWaitPercent)
	parts = append(parts, iowaitStyle.Render(iowait))

	return strings.Join(parts, " ")
}