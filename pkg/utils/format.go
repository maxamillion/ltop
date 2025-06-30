package utils

import (
	"fmt"
	"time"
)

func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func FormatBytesPerSecond(bytesPerSec float64) string {
	return FormatBytes(uint64(bytesPerSec)) + "/s"
}

func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	return fmt.Sprintf("%dd %dh", days, hours)
}

func FormatPercent(value float64) string {
	return fmt.Sprintf("%.1f%%", value)
}

func FormatFloat(value float64, precision int) string {
	return fmt.Sprintf("%."+fmt.Sprintf("%d", precision)+"f", value)
}

func FormatHz(hz uint64) string {
	if hz < 1000 {
		return fmt.Sprintf("%d Hz", hz)
	}
	if hz < 1000000 {
		return fmt.Sprintf("%.1f KHz", float64(hz)/1000)
	}
	if hz < 1000000000 {
		return fmt.Sprintf("%.1f MHz", float64(hz)/1000000)
	}
	return fmt.Sprintf("%.1f GHz", float64(hz)/1000000000)
}

func FormatTemperature(celsius float64) string {
	return fmt.Sprintf("%.1f°C", celsius)
}

func FormatUptime(uptime time.Duration) string {
	totalSeconds := int(uptime.Seconds())
	days := totalSeconds / 86400
	hours := (totalSeconds % 86400) / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	if days > 0 {
		return fmt.Sprintf("%dd %02d:%02d:%02d", days, hours, minutes, seconds)
	}
	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func PadString(s string, length int, padChar rune) string {
	if len(s) >= length {
		return s
	}
	padding := length - len(s)
	padStr := ""
	for i := 0; i < padding; i++ {
		padStr += string(padChar)
	}
	return s + padStr
}

func FormatLoadAverage(load [3]float64) string {
	return fmt.Sprintf("%.2f %.2f %.2f", load[0], load[1], load[2])
}

func ColorByPercentage(percentage float64) string {
	if percentage < 50 {
		return "green"
	} else if percentage < 80 {
		return "yellow"
	}
	return "red"
}

func FormatProcessState(state string) string {
	switch state {
	case "R":
		return "Running"
	case "S":
		return "Sleeping"
	case "D":
		return "Disk Sleep"
	case "T":
		return "Stopped"
	case "t":
		return "Traced"
	case "Z":
		return "Zombie"
	case "X":
		return "Dead"
	case "x":
		return "Dead"
	case "K":
		return "Wakekill"
	case "W":
		return "Waking"
	case "P":
		return "Parked"
	case "I":
		return "Idle"
	default:
		return "Unknown"
	}
}

func FormatTime(t time.Time) string {
	return t.Format("15:04:05")
}

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}