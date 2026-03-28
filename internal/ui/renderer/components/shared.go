package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
)

func RenderRange(since, until time.Time) string {
	if since.IsZero() && until.IsZero() {
		return ""
	}
	if since.IsZero() {
		return fmt.Sprintf("until %s", until.Format("2006-01-02"))
	}
	if until.IsZero() {
		return fmt.Sprintf("since %s", since.Format("2006-01-02"))
	}
	return fmt.Sprintf("%s ~ %s", since.Format("2006-01-02"), until.Format("2006-01-02"))
}

func RenderSpecsLine(specs domain.HardwareSpecs) string {
	parts := make([]string, 0, 4)
	if specs.Device != "" {
		parts = append(parts, specs.Device)
	}
	if specs.CPU != "" {
		parts = append(parts, specs.CPU)
	}
	if specs.RAM != "" {
		parts = append(parts, specs.RAM)
	}
	if specs.Battery != "" {
		parts = append(parts, specs.Battery)
	}
	return strings.Join(parts, "  |  ")
}

func joinLines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = -d
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh%02dm", h, m)
}

func dischargeRate(s domain.Session) float64 {
	dur := s.End.Sub(s.Start).Hours()
	if dur <= 0 {
		return 0
	}
	drain := s.StartPct - s.EndPct
	if drain <= 0 {
		return 0
	}
	return drain / dur
}

func levelString(level domain.LoadLevel) string {
	switch level {
	case domain.LoadLevelLight:
		return "light"
	case domain.LoadLevelMedium:
		return "medium"
	case domain.LoadLevelHeavy:
		return "heavy"
	default:
		return "unknown"
	}
}
