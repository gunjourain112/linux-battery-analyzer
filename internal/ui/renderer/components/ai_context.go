package components

import (
	"fmt"
	"strings"

	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/i18n"
)

func AIContext(tr i18n.Translator, config domain.Config, specs domain.HardwareSpecs, sessions []domain.Session, battery []domain.BatteryPoint, processes []domain.ProcessUsage, impacts []domain.ProcessImpact, profile domain.DischargeProfile, thermal domain.ThermalStats) string {
	lines := make([]string, 0, 16)

	lines = append(lines, tr.Get(i18n.AIContextLead))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Period: %s ~ %s", config.Since.Format("2006-01-02"), config.Until.Format("2006-01-02")))
	lines = append(lines, fmt.Sprintf("Device: %s | OS: %s | CPU: %s | RAM: %s", blankDash(specs.Device), blankDash(specs.OS), blankDash(specs.CPU), blankDash(specs.RAM)))
	lines = append(lines, fmt.Sprintf("Battery: %s", blankDash(specs.Battery)))
	if cur, start, end, ok := batteryWindow(battery); ok {
		lines = append(lines, fmt.Sprintf("Observed battery: %.0f%% -> %.0f%% (%s)", start.Percentage, end.Percentage, formatDuration(end.Time.Sub(start.Time))))
		lines = append(lines, fmt.Sprintf("Current battery: %.0f%%", cur.Percentage))
	}
	lines = append(lines, fmt.Sprintf("Sessions: %d | Process samples: %d | Battery buckets: %d", len(sessions), len(processes), profile.TotalCount))
	if avg := averageSessionDischargeRate(sessions); avg > 0 {
		lines = append(lines, fmt.Sprintf("Average session drain: %.2f%%/h", avg))
	}
	if thermal.Count > 0 {
		lines = append(lines, fmt.Sprintf("Thermals: %d samples, min %dC, max %dC, avg %dC", thermal.Count, thermal.Min, thermal.Max, thermal.Avg))
	}
	if top := topProcessImpact(impacts); top != nil {
		lines = append(lines, fmt.Sprintf("Top drain process: %s (%.1fW)", top.Process.Name, top.DrainWatts))
	}
	if bucket := dominantLoadBucket(profile); bucket != nil {
		lines = append(lines, fmt.Sprintf("Dominant discharge bucket: %s (%.0f%%, %.1fW)", bucket.Label, bucket.Ratio, bucket.AvgWatts))
	}
	if len(processes) > 0 {
		lines = append(lines, "")
		lines = append(lines, "Top processes:")
		for _, p := range topProcessLines(processes, 5) {
			lines = append(lines, "  - "+p)
		}
	}
	if len(profile.Buckets) > 0 {
		lines = append(lines, "")
		lines = append(lines, "Discharge profile:")
		for _, b := range profile.Buckets {
			if b.Count == 0 {
				continue
			}
			lines = append(lines, fmt.Sprintf("  - %s: %.0f%%, %.1fW, %s", b.Label, b.Ratio, b.AvgWatts, formatDurationFromHours(b.EstHours)))
		}
	}
	lines = append(lines, "")
	lines = append(lines, tr.Get(i18n.AIContextAsk))

	return strings.TrimRight(strings.Join(lines, "\n"), "\n")
}

func blankDash(v string) string {
	if strings.TrimSpace(v) == "" {
		return "--"
	}
	return v
}

func batteryWindow(points []domain.BatteryPoint) (current domain.BatteryPoint, start domain.BatteryPoint, end domain.BatteryPoint, ok bool) {
	if len(points) == 0 {
		return domain.BatteryPoint{}, domain.BatteryPoint{}, domain.BatteryPoint{}, false
	}

	current = points[0]
	for _, p := range points[1:] {
		if p.Time.After(current.Time) {
			current = p
		}
	}

	start = current
	end = current
	found := false
	for _, p := range points {
		if isObservedBatteryPoint(p) {
			if !found {
				start = p
				found = true
			}
			end = p
		}
	}
	if !found {
		start = points[0]
		end = points[len(points)-1]
		if end.Time.Before(start.Time) {
			start, end = end, start
		}
		if end.Time.Sub(start.Time) <= 0 {
			return domain.BatteryPoint{}, domain.BatteryPoint{}, domain.BatteryPoint{}, false
		}
		return current, start, end, true
	}
	if end.Time.Before(start.Time) {
		start, end = end, start
	}
	if end.Time.Sub(start.Time) <= 0 {
		return domain.BatteryPoint{}, domain.BatteryPoint{}, domain.BatteryPoint{}, false
	}
	return current, start, end, true
}

func isObservedBatteryPoint(p domain.BatteryPoint) bool {
	if strings.ToLower(strings.TrimSpace(p.State)) == "fully-charged" {
		return p.Percentage < 99.5
	}
	return true
}

func topProcessLines(processes []domain.ProcessUsage, limit int) []string {
	aggs := aggregateProcesses(processes)
	if len(aggs) == 0 {
		return nil
	}

	sortProcessAggregates(aggs)
	if limit > len(aggs) {
		limit = len(aggs)
	}

	lines := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		p := aggs[i]
		lines = append(lines, fmt.Sprintf("%s (%d samples, %.0fs CPU, %.1f max MEM)", p.Name, p.Count, p.CPUTime, p.MemPeak))
	}
	return lines
}

func sortProcessAggregates(aggs []processAggregate) {
	for i := 1; i < len(aggs); i++ {
		j := i
		for j > 0 {
			left := aggs[j-1]
			right := aggs[j]
			if left.Count > right.Count || (left.Count == right.Count && left.CPUTime >= right.CPUTime) {
				break
			}
			aggs[j-1], aggs[j] = aggs[j], aggs[j-1]
			j--
		}
	}
}
