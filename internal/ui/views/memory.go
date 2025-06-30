package views

import (
	"fmt"
	"strings"

	"github.com/admiller/ltop/internal/models"
	"github.com/admiller/ltop/internal/ui/components"
	"github.com/admiller/ltop/internal/ui/styles"
	"github.com/admiller/ltop/pkg/utils"
)

type MemoryView struct {
	memoryGauge *components.Gauge
	swapGauge   *components.Gauge
	multiGauge  *components.MultiGauge
}

func NewMemoryView() *MemoryView {
	return &MemoryView{
		memoryGauge: components.NewGauge(40),
		swapGauge:   components.NewGauge(40),
		multiGauge:  components.NewMultiGauge(60),
	}
}

func (mv *MemoryView) Render(snapshot *models.MetricsSnapshot, width, height int) string {
	if snapshot == nil {
		return "No data available"
	}

	var sections []string

	sections = append(sections, mv.renderMemoryOverview(snapshot))
	sections = append(sections, mv.renderMemoryBreakdown(snapshot))
	sections = append(sections, mv.renderSwapUsage(snapshot))

	content := strings.Join(sections, "\n\n")
	return styles.Panel().Width(width).Render(content)
}

func (mv *MemoryView) renderMemoryOverview(snapshot *models.MetricsSnapshot) string {
	var info []string
	info = append(info, styles.Title().Render("Memory Usage"))

	memoryGauge := mv.memoryGauge.Render(snapshot.Memory.UsedPercent, "Memory")
	info = append(info, memoryGauge)

	detail := fmt.Sprintf("  %s used / %s total (%s available)",
		utils.FormatBytes(snapshot.Memory.Used),
		utils.FormatBytes(snapshot.Memory.Total),
		utils.FormatBytes(snapshot.Memory.Available))
	info = append(info, styles.Muted().Render(detail))

	return strings.Join(info, "\n")
}

func (mv *MemoryView) renderMemoryBreakdown(snapshot *models.MetricsSnapshot) string {
	var breakdown []string
	breakdown = append(breakdown, styles.Title().Render("Memory Breakdown"))

	mv.multiGauge.Clear()
	mv.multiGauge.AddSegment(snapshot.Memory.Used, "Used", styles.Error())
	mv.multiGauge.AddSegment(snapshot.Memory.Cached, "Cached", styles.Warning())
	mv.multiGauge.AddSegment(snapshot.Memory.Buffers, "Buffers", styles.Info())

	gauge := mv.multiGauge.Render(snapshot.Memory.Total)
	breakdown = append(breakdown, gauge)

	breakdown = append(breakdown, "")
	breakdown = append(breakdown, fmt.Sprintf("Used:    %s", utils.FormatBytes(snapshot.Memory.Used)))
	breakdown = append(breakdown, fmt.Sprintf("Cached:  %s", utils.FormatBytes(snapshot.Memory.Cached)))
	breakdown = append(breakdown, fmt.Sprintf("Buffers: %s", utils.FormatBytes(snapshot.Memory.Buffers)))
	breakdown = append(breakdown, fmt.Sprintf("Free:    %s", utils.FormatBytes(snapshot.Memory.Free)))

	if snapshot.Memory.Shared > 0 {
		breakdown = append(breakdown, fmt.Sprintf("Shared:  %s", utils.FormatBytes(snapshot.Memory.Shared)))
	}

	return strings.Join(breakdown, "\n")
}

func (mv *MemoryView) renderSwapUsage(snapshot *models.MetricsSnapshot) string {
	var swap []string
	swap = append(swap, styles.Title().Render("Swap Usage"))

	if snapshot.Memory.Swap.Total == 0 {
		swap = append(swap, styles.Muted().Render("No swap configured"))
		return strings.Join(swap, "\n")
	}

	swapGauge := mv.swapGauge.Render(snapshot.Memory.Swap.UsedPercent, "Swap")
	swap = append(swap, swapGauge)

	detail := fmt.Sprintf("  %s used / %s total",
		utils.FormatBytes(snapshot.Memory.Swap.Used),
		utils.FormatBytes(snapshot.Memory.Swap.Total))
	swap = append(swap, styles.Muted().Render(detail))

	return strings.Join(swap, "\n")
}