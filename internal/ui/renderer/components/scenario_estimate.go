package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/i18n"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func ScenarioEstimate(profile domain.DischargeProfile, specs domain.HardwareSpecs, points []domain.BatteryPoint, tr i18n.Translator) string {
	current, ok := latestBatteryPoint(points)
	if !ok || profile.TotalCount == 0 {
		return theme.Default.Subtle().Render(tr.Get(i18n.NoScenarioEstimateData))
	}

	_, designWh, _, ok := parseBatterySpec(specs.Battery)
	if !ok || designWh <= 0 {
		return theme.Default.Subtle().Render(tr.Get(i18n.NoScenarioEstimateData))
	}

	currentWh := (current.Percentage / 100.0) * designWh
	fullWh := designWh

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers(
			tr.Get(i18n.BucketHeader),
			tr.Get(i18n.RatioHeader),
			tr.Get(i18n.AvgWHeader),
			tr.Get(i18n.CurrentHeader),
			tr.Get(i18n.FullHeader),
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

		tbl.Row(
			b.Label,
			fmt.Sprintf("%.0f%%", b.Ratio),
			fmt.Sprintf("%.1fW", b.AvgWatts),
			formatDurationFromHours(currentWh/b.AvgWatts),
			formatDurationFromHours(fullWh/b.AvgWatts),
		)
		rows++
	}

	if rows == 0 {
		return theme.Default.Subtle().Render(tr.Get(i18n.NoScenarioEstimateData))
	}

	out := tbl.Render()
	out += "\n" + theme.Default.Subtle().Render(tr.Get(i18n.ScenarioEstimateNote))
	return out
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
