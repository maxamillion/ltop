package components

import (
	"strings"

	"github.com/admiller/ltop/internal/ui/styles"
)

type ConfirmDialog struct {
	Title       string
	Message     string
	Visible     bool
	Width       int
	Height      int
	ConfirmText string
	CancelText  string
	Selected    int
}

func NewConfirmDialog(title, message string) *ConfirmDialog {
	return &ConfirmDialog{
		Title:       title,
		Message:     message,
		Visible:     false,
		Width:       50,
		Height:      8,
		ConfirmText: "Yes",
		CancelText:  "No",
		Selected:    1,
	}
}

func (cd *ConfirmDialog) Show() {
	cd.Visible = true
	cd.Selected = 1
}

func (cd *ConfirmDialog) Hide() {
	cd.Visible = false
}

func (cd *ConfirmDialog) IsVisible() bool {
	return cd.Visible
}

func (cd *ConfirmDialog) MoveLeft() {
	if cd.Selected > 0 {
		cd.Selected--
	}
}

func (cd *ConfirmDialog) MoveRight() {
	if cd.Selected < 1 {
		cd.Selected++
	}
}

func (cd *ConfirmDialog) IsConfirmSelected() bool {
	return cd.Selected == 0
}

func (cd *ConfirmDialog) Render() string {
	if !cd.Visible {
		return ""
	}

	var content []string

	content = append(content, styles.Title().Render(cd.Title))
	content = append(content, "")

	messageLines := cd.wrapText(cd.Message, cd.Width-4)
	content = append(content, messageLines...)

	content = append(content, "")

	buttons := cd.renderButtons()
	content = append(content, buttons)

	dialogContent := strings.Join(content, "\n")

	return styles.Border().
		Width(cd.Width).
		Height(cd.Height).
		Render(dialogContent)
}

func (cd *ConfirmDialog) wrapText(text string, width int) []string {
	if len(text) <= width {
		return []string{text}
	}

	var lines []string
	words := strings.Fields(text)
	currentLine := ""

	for _, word := range words {
		if len(currentLine)+len(word)+1 <= width {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

func (cd *ConfirmDialog) renderButtons() string {
	confirmStyle := styles.TableRow()
	cancelStyle := styles.TableRow()

	if cd.Selected == 0 {
		confirmStyle = styles.TableRowSelected()
	} else {
		cancelStyle = styles.TableRowSelected()
	}

	confirmButton := confirmStyle.Render(" " + cd.ConfirmText + " ")
	cancelButton := cancelStyle.Render(" " + cd.CancelText + " ")

	return "  " + confirmButton + "  " + cancelButton
}

type InputDialog struct {
	Title       string
	Message     string
	Input       *TextInput
	Visible     bool
	Width       int
	Height      int
	ConfirmText string
	CancelText  string
}

func NewInputDialog(title, message, placeholder string) *InputDialog {
	return &InputDialog{
		Title:       title,
		Message:     message,
		Input:       NewTextInput(placeholder, 30),
		Visible:     false,
		Width:       50,
		Height:      10,
		ConfirmText: "OK",
		CancelText:  "Cancel",
	}
}

func (id *InputDialog) Show() {
	id.Visible = true
	id.Input.Focus()
	id.Input.Clear()
}

func (id *InputDialog) Hide() {
	id.Visible = false
	id.Input.Blur()
}

func (id *InputDialog) IsVisible() bool {
	return id.Visible
}

func (id *InputDialog) GetValue() string {
	return id.Input.GetValue()
}

func (id *InputDialog) HandleInput(ch rune) {
	id.Input.InsertChar(ch)
}

func (id *InputDialog) HandleBackspace() {
	id.Input.DeleteChar()
}

func (id *InputDialog) Render() string {
	if !id.Visible {
		return ""
	}

	var content []string

	content = append(content, styles.Title().Render(id.Title))
	content = append(content, "")

	messageLines := strings.Split(id.Message, "\n")
	content = append(content, messageLines...)

	content = append(content, "")
	content = append(content, "Input: "+id.Input.Render())
	content = append(content, "")
	content = append(content, styles.HelpText().Render("Press Enter to confirm, Esc to cancel"))

	dialogContent := strings.Join(content, "\n")

	return styles.Border().
		Width(id.Width).
		Height(id.Height).
		Render(dialogContent)
}
