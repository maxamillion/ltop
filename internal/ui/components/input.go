package components

import (
	"strings"

	"github.com/admiller/ltop/internal/ui/styles"
)

type TextInput struct {
	Value       string
	Placeholder string
	Width       int
	Focused     bool
	CursorPos   int
}

func NewTextInput(placeholder string, width int) *TextInput {
	return &TextInput{
		Placeholder: placeholder,
		Width:       width,
		Focused:     false,
		CursorPos:   0,
	}
}

func (ti *TextInput) SetValue(value string) {
	ti.Value = value
	if ti.CursorPos > len(value) {
		ti.CursorPos = len(value)
	}
}

func (ti *TextInput) Focus() {
	ti.Focused = true
}

func (ti *TextInput) Blur() {
	ti.Focused = false
}

func (ti *TextInput) InsertChar(ch rune) {
	if ti.CursorPos <= len(ti.Value) {
		ti.Value = ti.Value[:ti.CursorPos] + string(ch) + ti.Value[ti.CursorPos:]
		ti.CursorPos++
	}
}

func (ti *TextInput) DeleteChar() {
	if ti.CursorPos > 0 && len(ti.Value) > 0 {
		ti.Value = ti.Value[:ti.CursorPos-1] + ti.Value[ti.CursorPos:]
		ti.CursorPos--
	}
}

func (ti *TextInput) MoveCursorLeft() {
	if ti.CursorPos > 0 {
		ti.CursorPos--
	}
}

func (ti *TextInput) MoveCursorRight() {
	if ti.CursorPos < len(ti.Value) {
		ti.CursorPos++
	}
}

func (ti *TextInput) MoveCursorToStart() {
	ti.CursorPos = 0
}

func (ti *TextInput) MoveCursorToEnd() {
	ti.CursorPos = len(ti.Value)
}

func (ti *TextInput) Clear() {
	ti.Value = ""
	ti.CursorPos = 0
}

func (ti *TextInput) Render() string {
	text := ti.Value
	if text == "" && !ti.Focused {
		text = ti.Placeholder
		return styles.Muted().Render(text)
	}

	if ti.Focused {
		if ti.CursorPos < len(text) {
			before := text[:ti.CursorPos]
			cursor := text[ti.CursorPos : ti.CursorPos+1]
			after := text[ti.CursorPos+1:]
			text = before + styles.TableRowSelected().Render(cursor) + after
		} else {
			text = text + styles.TableRowSelected().Render(" ")
		}
	}

	if len(text) > ti.Width {
		if ti.CursorPos > ti.Width-1 {
			start := ti.CursorPos - ti.Width + 1
			text = text[start:]
		} else {
			text = text[:ti.Width]
		}
	}

	padding := ti.Width - len(strings.ReplaceAll(text, "\x1b", ""))
	if padding > 0 {
		text += strings.Repeat(" ", padding)
	}

	style := styles.TableRow()
	if ti.Focused {
		style = styles.Border()
	}

	return style.Render(text)
}

func (ti *TextInput) GetValue() string {
	return ti.Value
}

func (ti *TextInput) IsEmpty() bool {
	return ti.Value == ""
}