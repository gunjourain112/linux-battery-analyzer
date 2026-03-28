package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/i18n"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func HeaderSummary(tr i18n.Translator, specs domain.HardwareSpecs, sessions []domain.Session, battery []domain.BatteryPoint, profile domain.DischargeProfile, thermal domain.ThermalStats) string {
	current, ok := latestBatteryPoint(battery)
	if !ok {
		return ""
	}

	_, designWh, _, specOK := parseBatterySpec(specs.Battery)
	if !specOK || designWh <= 0 {
		designWh = 0
	}

	currentWatts := weightedAverageWatts(profile)
	currentEstimate := estimateBatteryHours(current.Percentage, designWh, currentWatts)
	rangeIdle, rangeHeavy := estimateBatteryRange(current.Percentage, designWh, profile)
	actualUse := totalSessionDuration(sessions)
	worst := topWorstSession(sessions, designWh)

	lines := []string{
		fmt.Sprintf("%s %.0f%% %s → %s (%.0f%% / %.1fW 기준)",
			tr.Get(i18n.CurrentStatusLabel),
			current.Percentage,
			batteryStateLabel(current.State),
			currentEstimate,
			current.Percentage,
			currentWatts,
		),
		fmt.Sprintf("%s: %s / %s", tr.Get(i18n.RangeSummaryLabel), rangeIdle, rangeHeavy),
		fmt.Sprintf("%s %s", tr.Get(i18n.ActualUseLabel), formatDurationCard(actualUse)),
	}

	if worst != "" {
		lines = append(lines, fmt.Sprintf("%s %s", tr.Get(i18n.TopDrainSessionLabel), worst))
	}
	if thermal.Count > 0 {
		lines = append(lines, fmt.Sprintf("%s: %d°C", tr.Get(i18n.PeakTempLabel), thermal.Max))
	}

	body := strings.Join(lines, "\n")
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(theme.Default.Gray).
		Padding(0, 1)
	return box.Render(body)
}

func estimateBatteryHours(pct float64, capacityWh float64, watts float64) string {
	if capacityWh <= 0 || watts <= 0 || pct <= 0 {
		return "--"
	}
	remainingWh := capacityWh * pct / 100
	return fmt.Sprintf("약 %s", formatDurationFromHours(remainingWh/watts))
}

func estimateBatteryRange(currentPct float64, capacityWh float64, profile domain.DischargeProfile) (string, string) {
	if capacityWh <= 0 {
		return "--", "--"
	}

	var low *domain.LoadBucket
	var high *domain.LoadBucket
	for i := range profile.Buckets {
		b := &profile.Buckets[i]
		if b.Count == 0 || b.AvgWatts <= 0 {
			continue
		}
		if low == nil || b.AvgWatts < low.AvgWatts {
			low = b
		}
		if high == nil || b.AvgWatts > high.AvgWatts {
			high = b
		}
	}

	currentWh := capacityWh * currentPct / 100
	if low == nil || high == nil {
		return "--", "--"
	}

	return fmt.Sprintf("%s ~%s", low.Label, formatDurationFromHours(currentWh/low.AvgWatts)),
		fmt.Sprintf("%s ~%s", high.Label, formatDurationFromHours(currentWh/high.AvgWatts))
}

func totalSessionDuration(sessions []domain.Session) time.Duration {
	var total time.Duration
	for _, s := range sessions {
		if s.End.After(s.Start) {
			total += s.End.Sub(s.Start)
		}
	}
	return total
}

func topWorstSession(sessions []domain.Session, capacityWh float64) string {
	var worst domain.Session
	var hasWorst bool
	var worstRate float64
	for _, s := range sessions {
		rate := dischargeRate(s)
		if rate <= 0 {
			continue
		}
		if !hasWorst || rate > worstRate {
			worst = s
			worstRate = rate
			hasWorst = true
		}
	}
	if !hasWorst {
		return ""
	}

	avgWatts := "--"
	if capacityWh > 0 {
		drain := worst.StartPct - worst.EndPct
		dur := worst.End.Sub(worst.Start).Hours()
		if dur > 0 && drain > 0 {
			avgWatts = fmt.Sprintf("%.1fW", (capacityWh*drain/100)/dur)
		}
	}

	return fmt.Sprintf("%s ~ %s (%s / -%.0f%%p / avg %s)",
		worst.Start.Format("01/02 15:04"),
		worst.End.Format("15:04"),
		formatDuration(worst.End.Sub(worst.Start)),
		worst.StartPct-worst.EndPct,
		avgWatts,
	)
}

func formatDurationCard(d time.Duration) string {
	if d < 0 {
		d = -d
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh %02dm", h, m)
}

func batteryStateLabel(state string) string {
	switch strings.ToLower(strings.TrimSpace(state)) {
	case "charging":
		return "(charging)"
	case "fully-charged":
		return "(fully-charged)"
	case "discharging":
		return "(discharging)"
	default:
		return "(--)"
	}
}
