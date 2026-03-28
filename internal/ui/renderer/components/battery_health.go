package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/i18n"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func BatteryHealth(specs domain.HardwareSpecs, points []domain.BatteryPoint, tr i18n.Translator) string {
	rows := make([][]string, 0, 6)
	if fullWh, designWh, healthPct, ok := parseBatterySpec(specs.Battery); ok {
		rows = append(rows, []string{tr.Get(i18n.BatteryHeader), specs.Battery})
		rows = append(rows,
			[]string{tr.Get(i18n.DesignCapacityHeader), fmt.Sprintf("%.1f Wh", designWh)},
			[]string{tr.Get(i18n.CurrentCapacityHeader), fmt.Sprintf("%.1f Wh", fullWh)},
		)
		if healthPct > 0 {
			rows = append(rows, []string{tr.Get(i18n.HealthHeader), fmt.Sprintf("%.1f%%", healthPct)})
		}
	}

	if start, end, ok := batteryRange(points); ok {
		rows = append(rows,
			[]string{tr.Get(i18n.SamplesHeader), fmt.Sprintf("%d", len(points))},
			[]string{tr.Get(i18n.StartHeader), fmt.Sprintf("%.0f%%", start.Percentage)},
			[]string{tr.Get(i18n.EndHeader), fmt.Sprintf("%.0f%%", end.Percentage)},
			[]string{tr.Get(i18n.DurationHeader), formatDuration(end.Time.Sub(start.Time))},
			[]string{tr.Get(i18n.DrainHeader), fmt.Sprintf("%+.1f%%", end.Percentage-start.Percentage)},
			[]string{tr.Get(i18n.RateHeader), fmt.Sprintf("%.2f%%/h", batteryDrainRate(start, end))},
		)
	}

	if len(rows) == 0 {
		return theme.Default.Subtle().Render(tr.Get(i18n.NoBatteryHealthData))
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers(tr.Get(i18n.MetricHeader), tr.Get(i18n.ValueHeader)).
		StyleFunc(func(r, c int) lipgloss.Style {
			if r == -1 {
				return theme.Default.Header()
			}
			if c == 0 {
				return theme.Default.Label()
			}
			return theme.Default.Value()
		})

	for _, row := range rows {
		tbl.Row(row...)
	}

	return tbl.Render()
}

func batteryRange(points []domain.BatteryPoint) (domain.BatteryPoint, domain.BatteryPoint, bool) {
	if len(points) < 2 {
		return domain.BatteryPoint{}, domain.BatteryPoint{}, false
	}

	start := points[0]
	end := points[len(points)-1]
	if end.Time.Before(start.Time) {
		start, end = end, start
	}
	if end.Time.Sub(start.Time) <= 0 {
		return domain.BatteryPoint{}, domain.BatteryPoint{}, false
	}
	return start, end, true
}

func batteryDrainRate(start, end domain.BatteryPoint) float64 {
	hours := end.Time.Sub(start.Time).Hours()
	if hours <= 0 {
		return 0
	}
	return (start.Percentage - end.Percentage) / hours
}
