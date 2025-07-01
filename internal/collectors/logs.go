package collectors

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/admiller/ltop/internal/models"
)

type LogCollector struct {
	sources       []string
	maxEntries    int
	lastEntries   map[string]time.Time
	logLevelRegex *regexp.Regexp
}

func NewLogCollector(sources []string, maxEntries int) *LogCollector {
	if len(sources) == 0 {
		sources = []string{"/var/log/syslog", "journalctl"}
	}
	if maxEntries == 0 {
		maxEntries = 50
	}

	logLevelRegex := regexp.MustCompile(`(?i)\b(emerg|alert|crit|err|error|warn|warning|notice|info|debug)\b`)

	return &LogCollector{
		sources:       sources,
		maxEntries:    maxEntries,
		lastEntries:   make(map[string]time.Time),
		logLevelRegex: logLevelRegex,
	}
}

func (l *LogCollector) Collect() (*models.LogMetrics, error) {
	metrics := &models.LogMetrics{
		Timestamp: time.Now(),
		Sources:   l.sources,
	}

	allEntries := make([]models.LogEntry, 0)

	for _, source := range l.sources {
		entries, err := l.collectFromSource(source)
		if err != nil {
			continue
		}
		allEntries = append(allEntries, entries...)
	}

	l.sortAndLimitEntries(&allEntries)
	l.categorizeEntries(allEntries, metrics)

	metrics.Entries = allEntries
	return metrics, nil
}

func (l *LogCollector) collectFromSource(source string) ([]models.LogEntry, error) {
	if source == "journalctl" {
		return l.collectFromJournalctl()
	}
	return l.collectFromFile(source)
}

func (l *LogCollector) collectFromFile(filename string) ([]models.LogEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	entries := make([]models.LogEntry, 0)
	scanner := bufio.NewScanner(file)

	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	start := len(lines) - l.maxEntries
	if start < 0 {
		start = 0
	}

	for i := start; i < len(lines); i++ {
		entry := l.parseLogLine(lines[i], filename)
		if !entry.Timestamp.IsZero() {
			entries = append(entries, entry)
		}
	}

	return entries, scanner.Err()
}

func (l *LogCollector) collectFromJournalctl() ([]models.LogEntry, error) {
	cmd := exec.Command("journalctl", "-n", fmt.Sprintf("%d", l.maxEntries), "--no-pager", "-o", "short-iso")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	entries := make([]models.LogEntry, 0)
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		entry := l.parseJournalLine(line)
		if !entry.Timestamp.IsZero() {
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

func (l *LogCollector) parseLogLine(line, source string) models.LogEntry {
	entry := models.LogEntry{
		Source:  source,
		Message: line,
	}

	entry.Timestamp = l.extractTimestamp(line)
	entry.Level = l.extractLogLevel(line)
	entry.Service = l.extractService(line)

	return entry
}

func (l *LogCollector) parseJournalLine(line string) models.LogEntry {
	entry := models.LogEntry{
		Source: "journalctl",
	}

	parts := strings.SplitN(line, " ", 4)
	if len(parts) < 4 {
		entry.Message = line
		return entry
	}

	// The timestamp is only the first part, parts[1] is the hostname
	timestampStr := parts[0]
	if timestamp, err := time.Parse("2006-01-02T15:04:05-07:00", timestampStr); err == nil {
		entry.Timestamp = timestamp
	} else if timestamp, err := time.Parse("2006-01-02T15:04:05+00:00", timestampStr); err == nil {
		entry.Timestamp = timestamp
	} else if timestamp, err := time.Parse("2006-01-02T15:04:05-0700", timestampStr); err == nil {
		entry.Timestamp = timestamp
	} else if timestamp, err := time.Parse("2006-01-02T15:04:05+0000", timestampStr); err == nil {
		entry.Timestamp = timestamp
	}

	entry.Service = parts[2]
	entry.Message = parts[3]
	entry.Level = l.extractLogLevel(entry.Message)

	return entry
}

func (l *LogCollector) extractTimestamp(line string) time.Time {
	now := time.Now()
	currentYear := now.Year()

	timestampFormats := []string{
		"Jan 2 15:04:05",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05-0700",
		"2006-01-02 15:04:05",
		"Jan  2 15:04:05",
	}

	for _, format := range timestampFormats {
		if timestamp, err := time.Parse(format, line[:min(len(line), len(format)+5)]); err == nil {
			if timestamp.Year() == 0 {
				timestamp = timestamp.AddDate(currentYear, 0, 0)
			}
			return timestamp
		}
	}

	return time.Time{}
}

func (l *LogCollector) extractLogLevel(message string) string {
	message = strings.ToLower(message)

	matches := l.logLevelRegex.FindStringSubmatch(message)
	if len(matches) > 1 {
		level := strings.ToLower(matches[1])
		switch level {
		case "error", "err":
			return "err"
		case "warning", "warn":
			return "warn"
		default:
			return level
		}
	}

	if strings.Contains(message, "error") || strings.Contains(message, "fail") {
		return "err"
	}
	if strings.Contains(message, "warn") {
		return "warn"
	}
	if strings.Contains(message, "info") {
		return "info"
	}

	return "info"
}

func (l *LogCollector) extractService(line string) string {
	if strings.Contains(line, "systemd") {
		return "systemd"
	}
	if strings.Contains(line, "kernel") {
		return "kernel"
	}
	if strings.Contains(line, "sshd") {
		return "sshd"
	}
	if strings.Contains(line, "NetworkManager") {
		return "NetworkManager"
	}

	words := strings.Fields(line)
	if len(words) > 3 {
		for _, word := range words[3:6] {
			if strings.HasSuffix(word, ":") {
				return strings.TrimSuffix(word, ":")
			}
			if strings.Contains(word, "[") && strings.Contains(word, "]") {
				return word
			}
		}
	}

	return "system"
}

func (l *LogCollector) sortAndLimitEntries(entries *[]models.LogEntry) {
	if len(*entries) <= l.maxEntries {
		return
	}

	sorted := make([]models.LogEntry, len(*entries))
	copy(sorted, *entries)

	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i].Timestamp.Before(sorted[j].Timestamp) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	if len(sorted) > l.maxEntries {
		sorted = sorted[:l.maxEntries]
	}

	*entries = sorted
}

func (l *LogCollector) categorizeEntries(entries []models.LogEntry, metrics *models.LogMetrics) {
	for _, entry := range entries {
		switch entry.Level {
		case "emerg", "alert", "crit", "err", "error":
			metrics.ErrorCount++
		case "warn", "warning":
			metrics.WarnCount++
		default:
			metrics.InfoCount++
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
