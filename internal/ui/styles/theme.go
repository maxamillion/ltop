package styles

import (
	"github.com/admiller/ltop/internal/models"
	"github.com/charmbracelet/lipgloss"
)

var (
	DefaultTheme = models.DarkTheme()
)

func Base() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(DefaultTheme.Foreground)).
		Background(lipgloss.Color(DefaultTheme.Background))
}

func Header() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(DefaultTheme.Primary)).
		Background(lipgloss.Color(DefaultTheme.Background)).
		Padding(0, 1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(DefaultTheme.Primary))
}

func StatusBar() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(DefaultTheme.Background)).
		Background(lipgloss.Color(DefaultTheme.Primary)).
		Padding(0, 1)
}

func Title() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(DefaultTheme.Primary)).
		MarginBottom(1)
}

func TableHeader() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(DefaultTheme.Background)).
		Background(lipgloss.Color(DefaultTheme.Primary)).
		Padding(0, 1)
}

func TableRow() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(DefaultTheme.Foreground)).
		Padding(0, 1)
}

func TableRowSelected() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(DefaultTheme.Background)).
		Background(lipgloss.Color(DefaultTheme.Secondary)).
		Padding(0, 1)
}

func Gauge() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(DefaultTheme.Primary))
}

func GaugeWarning() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(DefaultTheme.Warning))
}

func GaugeCritical() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(DefaultTheme.Error))
}

func Success() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(DefaultTheme.Success))
}

func Warning() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(DefaultTheme.Warning))
}

func Error() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(DefaultTheme.Error))
}

func Info() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(DefaultTheme.Info))
}

func Muted() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(DefaultTheme.Muted))
}

func Border() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(DefaultTheme.Primary)).
		Padding(1)
}

func Panel() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(DefaultTheme.Secondary)).
		Padding(1)
}

func HelpText() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(DefaultTheme.Muted)).
		Italic(true)
}

func NavigationTab() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(DefaultTheme.Muted)).
		Padding(0, 2)
}

func NavigationTabActive() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(DefaultTheme.Primary)).
		Background(lipgloss.Color(DefaultTheme.Secondary)).
		Padding(0, 2).
		Bold(true)
}

func PercentageColor(percentage float64) lipgloss.Style {
	if percentage < 50 {
		return Success()
	} else if percentage < 80 {
		return Warning()
	}
	return Error()
}

func ProcessStateColor(state string) lipgloss.Style {
	switch state {
	case "R":
		return Success()
	case "S", "I":
		return Info()
	case "D":
		return Warning()
	case "T":
		return Warning()
	case "Z":
		return Error()
	default:
		return Muted()
	}
}

func LogLevelColor(level string) lipgloss.Style {
	switch level {
	case "emerg", "alert", "crit", "err", "error":
		return Error()
	case "warn", "warning":
		return Warning()
	case "info":
		return Info()
	case "debug":
		return Muted()
	default:
		return Base()
	}
}

func SetTheme(theme models.ColorTheme) {
	DefaultTheme = theme
}
