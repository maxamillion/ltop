package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/admiller/ltop/internal/ui/styles"
)

type Gauge struct {
	Width      int
	ShowValue  bool
	ShowLabel  bool
	BarChar    string
	EmptyChar  string
}

func NewGauge(width int) *Gauge {
	return &Gauge{
		Width:     width,
		ShowValue: true,
		ShowLabel: true,
		BarChar:   "█",
		EmptyChar: "░",
	}
}

func (g *Gauge) Render(value float64, label string) string {
	if value < 0 {
		value = 0
	}
	if value > 100 {
		value = 100
	}

	filled := int((value / 100.0) * float64(g.Width))
	empty := g.Width - filled

	bar := strings.Repeat(g.BarChar, filled) + strings.Repeat(g.EmptyChar, empty)

	style := styles.PercentageColor(value)
	bar = style.Render(bar)

	result := bar

	if g.ShowValue {
		valueStr := fmt.Sprintf(" %5.1f%%", value)
		result += valueStr
	}

	if g.ShowLabel && label != "" {
		result = label + ": " + result
	}

	return result
}

func (g *Gauge) RenderBytes(used, total uint64, label string) string {
	if total == 0 {
		return g.Render(0, label)
	}
	
	percentage := float64(used) / float64(total) * 100.0
	return g.Render(percentage, label)
}

func (g *Gauge) RenderWithColors(value float64, label string, lowColor, midColor, highColor lipgloss.Style) string {
	if value < 0 {
		value = 0
	}
	if value > 100 {
		value = 100
	}

	filled := int((value / 100.0) * float64(g.Width))
	empty := g.Width - filled

	bar := strings.Repeat(g.BarChar, filled) + strings.Repeat(g.EmptyChar, empty)

	var style lipgloss.Style
	if value < 50 {
		style = lowColor
	} else if value < 80 {
		style = midColor
	} else {
		style = highColor
	}

	bar = style.Render(bar)

	result := bar

	if g.ShowValue {
		valueStr := fmt.Sprintf(" %5.1f%%", value)
		result += valueStr
	}

	if g.ShowLabel && label != "" {
		result = label + ": " + result
	}

	return result
}

type MultiGauge struct {
	Width     int
	ShowTotal bool
	Segments  []GaugeSegment
}

type GaugeSegment struct {
	Value uint64
	Label string
	Color lipgloss.Style
}

func NewMultiGauge(width int) *MultiGauge {
	return &MultiGauge{
		Width:     width,
		ShowTotal: true,
		Segments:  make([]GaugeSegment, 0),
	}
}

func (mg *MultiGauge) AddSegment(value uint64, label string, color lipgloss.Style) {
	mg.Segments = append(mg.Segments, GaugeSegment{
		Value: value,
		Label: label,
		Color: color,
	})
}

func (mg *MultiGauge) Render(total uint64) string {
	if total == 0 {
		return strings.Repeat("░", mg.Width)
	}

	var result strings.Builder
	remaining := mg.Width
	
	for _, segment := range mg.Segments {
		if remaining <= 0 {
			break
		}
		
		segmentWidth := int(float64(segment.Value) / float64(total) * float64(mg.Width))
		if segmentWidth > remaining {
			segmentWidth = remaining
		}
		
		if segmentWidth > 0 {
			segmentBar := strings.Repeat("█", segmentWidth)
			result.WriteString(segment.Color.Render(segmentBar))
			remaining -= segmentWidth
		}
	}
	
	if remaining > 0 {
		result.WriteString(strings.Repeat("░", remaining))
	}
	
	return result.String()
}

func (mg *MultiGauge) Clear() {
	mg.Segments = mg.Segments[:0]
}