package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/admiller/ltop/internal/ui/styles"
	"github.com/admiller/ltop/pkg/utils"
)

type Table struct {
	Headers    []string
	Rows       [][]string
	Selected   int
	Width      int
	Height     int
	Scrollable bool
	ScrollTop  int
}

func NewTable(headers []string) *Table {
	return &Table{
		Headers:    headers,
		Rows:       make([][]string, 0),
		Selected:   0,
		Width:      80,
		Height:     20,
		Scrollable: true,
		ScrollTop:  0,
	}
}

func (t *Table) AddRow(row []string) {
	if len(row) != len(t.Headers) {
		for len(row) < len(t.Headers) {
			row = append(row, "")
		}
		row = row[:len(t.Headers)]
	}
	t.Rows = append(t.Rows, row)
}

func (t *Table) ClearRows() {
	t.Rows = t.Rows[:0]
	t.Selected = 0
	t.ScrollTop = 0
}

func (t *Table) SetSize(width, height int) {
	t.Width = width
	t.Height = height
}

func (t *Table) MoveUp() {
	if t.Selected > 0 {
		t.Selected--
		if t.Selected < t.ScrollTop {
			t.ScrollTop = t.Selected
		}
	}
}

func (t *Table) MoveDown() {
	if t.Selected < len(t.Rows)-1 {
		t.Selected++
		if t.Selected >= t.ScrollTop+t.Height-1 {
			t.ScrollTop = t.Selected - t.Height + 2
		}
	}
}

func (t *Table) PageUp() {
	t.Selected = utils.Max(0, t.Selected-t.Height+1)
	t.ScrollTop = utils.Max(0, t.ScrollTop-t.Height+1)
}

func (t *Table) PageDown() {
	maxSelected := len(t.Rows) - 1
	t.Selected = utils.Min(maxSelected, t.Selected+t.Height-1)
	t.ScrollTop = utils.Min(t.Selected, t.ScrollTop+t.Height-1)
}

func (t *Table) GetSelectedRow() []string {
	if t.Selected >= 0 && t.Selected < len(t.Rows) {
		return t.Rows[t.Selected]
	}
	return nil
}

func (t *Table) Render() string {
	if len(t.Headers) == 0 {
		return ""
	}

	colWidths := t.calculateColumnWidths()
	var result strings.Builder

	header := t.renderHeader(colWidths)
	result.WriteString(header)
	result.WriteString("\n")

	if len(t.Rows) == 0 {
		return result.String()
	}

	visibleRows := t.getVisibleRows()
	for i, row := range visibleRows {
		actualIndex := t.ScrollTop + i
		style := styles.TableRow()
		if actualIndex == t.Selected {
			style = styles.TableRowSelected()
		}
		
		renderedRow := t.renderRow(row, colWidths, style)
		result.WriteString(renderedRow)
		result.WriteString("\n")
	}

	return result.String()
}

func (t *Table) calculateColumnWidths() []int {
	if len(t.Headers) == 0 {
		return nil
	}

	availableWidth := t.Width - (len(t.Headers) - 1)
	colWidths := make([]int, len(t.Headers))
	
	totalContentWidth := 0
	for i, header := range t.Headers {
		maxWidth := len(header)
		for _, row := range t.Rows {
			if i < len(row) && len(row[i]) > maxWidth {
				maxWidth = len(row[i])
			}
		}
		colWidths[i] = maxWidth
		totalContentWidth += maxWidth
	}

	if totalContentWidth <= availableWidth {
		return colWidths
	}

	ratio := float64(availableWidth) / float64(totalContentWidth)
	for i := range colWidths {
		colWidths[i] = utils.Max(5, int(float64(colWidths[i])*ratio))
	}

	return colWidths
}

func (t *Table) renderHeader(colWidths []int) string {
	var parts []string
	for i, header := range t.Headers {
		width := colWidths[i]
		text := utils.TruncateString(header, width)
		text = utils.PadString(text, width, ' ')
		parts = append(parts, styles.TableHeader().Render(text))
	}
	return strings.Join(parts, " ")
}

func (t *Table) renderRow(row []string, colWidths []int, style lipgloss.Style) string {
	var parts []string
	for i, cell := range row {
		if i >= len(colWidths) {
			break
		}
		width := colWidths[i]
		text := utils.TruncateString(cell, width)
		text = utils.PadString(text, width, ' ')
		parts = append(parts, style.Render(text))
	}
	return strings.Join(parts, " ")
}

func (t *Table) getVisibleRows() [][]string {
	if len(t.Rows) == 0 {
		return nil
	}

	start := t.ScrollTop
	end := utils.Min(len(t.Rows), t.ScrollTop+t.Height-1)
	
	if start >= len(t.Rows) {
		start = len(t.Rows) - 1
		t.ScrollTop = start
	}
	
	if end > len(t.Rows) {
		end = len(t.Rows)
	}

	if start < 0 {
		start = 0
		t.ScrollTop = 0
	}

	return t.Rows[start:end]
}

func (t *Table) GetRowCount() int {
	return len(t.Rows)
}

func (t *Table) GetSelectedIndex() int {
	return t.Selected
}

func (t *Table) SetSelected(index int) {
	if index >= 0 && index < len(t.Rows) {
		t.Selected = index
		if t.Selected < t.ScrollTop {
			t.ScrollTop = t.Selected
		} else if t.Selected >= t.ScrollTop+t.Height-1 {
			t.ScrollTop = t.Selected - t.Height + 2
		}
	}
}

func (t *Table) RenderWithInfo() string {
	content := t.Render()
	
	if len(t.Rows) > 0 {
		info := fmt.Sprintf("Showing %d-%d of %d", 
			t.ScrollTop+1, 
			utils.Min(t.ScrollTop+t.Height-1, len(t.Rows)), 
			len(t.Rows))
		content += "\n" + styles.HelpText().Render(info)
	}
	
	return content
}