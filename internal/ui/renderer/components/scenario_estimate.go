package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/i18n"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func ScenarioEstimate(profile domain.DischargeProfile, points []domain.BatteryPoint, tr i18n.Translator) string {
	currentPct, drainRate, ok := observedBatteryTrend(points)
	if !ok || profile.TotalCount == 0 {
		return theme.Default.Subtle().Render(tr.Get(i18n.NoScenarioEstimateData))
	}

	observedWatts := weightedAverageWatts(profile)
	if observedWatts <= 0 {
		return theme.Default.Subtle().Render(tr.Get(i18n.NoScenarioEstimateData))
	}

	currentHours := currentPct / drainRate
	fullHours := 100.0 / drainRate

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers(
			tr.Get(i18n.BucketHeader),
			tr.Get(i18n.RatioHeader),
			tr.Get(i18n.AvgWHeader),
			tr.Get(i18n.DurationHeader),
			tr.Get(i18n.RateHeader),
		).
		StyleFunc(func(r, c int) lipgloss.Style {
			if r == -1 {
				return theme.Default.Header()
			}
			if c == 0 {
				return theme.Default.Label()
			}
			return theme.Default.Value().Align(lipgloss.Right)
		})

	rows := 0
	for _, b := range profile.Buckets {
		if b.Count == 0 || b.AvgWatts <= 0 {
			continue
		}

		scale := observedWatts / b.AvgWatts
		tbl.Row(
			b.Label,
			fmt.Sprintf("%.0f%%", b.Ratio),
			fmt.Sprintf("%.1fW", b.AvgWatts),
			formatDurationFromHours(currentHours*scale),
			formatDurationFromHours(fullHours*scale),
		)
		rows++
	}

	if rows == 0 {
		return theme.Default.Subtle().Render(tr.Get(i18n.NoScenarioEstimateData))
	}

	out := tbl.Render()
	out += "\n" + theme.Default.Subtle().Render(fmt.Sprintf("  %s: %.2f%%/h", tr.Get(i18n.RateHeader), drainRate))
	out += "\n" + theme.Default.Subtle().Render(tr.Get(i18n.ScenarioEstimateNote))
	return out
}

func observedBatteryTrend(points []domain.BatteryPoint) (currentPct float64, drainRate float64, ok bool) {
	if len(points) < 2 {
		return 0, 0, false
	}

	first := points[0]
	last := points[len(points)-1]
	if last.Time.Before(first.Time) {
		first, last = last, first
	}

	hours := last.Time.Sub(first.Time).Hours()
	if hours <= 0 {
		return 0, 0, false
	}

	drain := first.Percentage - last.Percentage
	if drain <= 0 {
		return 0, 0, false
	}

	return last.Percentage, drain / hours, true
}

func weightedAverageWatts(profile domain.DischargeProfile) float64 {
	var sum float64
	var weight float64
	for _, b := range profile.Buckets {
		if b.EstHours <= 0 || b.AvgWatts <= 0 {
			continue
		}
		sum += b.AvgWatts * b.EstHours
		weight += b.EstHours
	}
	if weight <= 0 {
		return 0
	}
	return sum / weight
}
