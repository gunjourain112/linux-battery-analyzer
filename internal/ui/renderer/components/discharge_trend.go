package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/i18n"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func DischargeTrend(tr i18n.Translator, rows []domain.DetailedTimelineRow) string {
	var filtered []domain.DetailedTimelineRow
	for _, row := range rows {
		if row.BatteryPct <= 0 && row.PowerWatts <= 0 && len(row.Events) == 0 && len(row.ActiveProcs) == 0 {
			continue
		}
		filtered = append(filtered, row)
	}
	if len(filtered) == 0 {
		return theme.Default.Subtle().Render(tr.Get(i18n.NoTimelineData))
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers(
			tr.Get(i18n.TimeHeader),
			tr.Get(i18n.BatteryHeader),
			tr.Get(i18n.PowerHeader),
			tr.Get(i18n.HeatHeader),
			tr.Get(i18n.StateHeader),
			tr.Get(i18n.ProcessHeader),
			tr.Get(i18n.EventHeader),
		).
		StyleFunc(func(r, c int) lipgloss.Style {
			if r == -1 {
				return theme.Default.Header()
			}
			if c == 0 {
				return theme.Default.Value()
			}
			return theme.Default.Value().Align(lipgloss.Right)
		})

	limit := len(filtered)
	if limit > 26 {
		limit = 26
	}

	for i := 0; i < limit; i++ {
		row := filtered[i]
		power := "--"
		if row.PowerWatts > 0 {
			power = fmt.Sprintf("%.1fW", row.PowerWatts)
		}
		temp := "--"
		if row.AvgTempC > 0 {
			temp = fmt.Sprintf("%.0f°C", row.AvgTempC)
		}
		processes := "--"
		if len(row.ActiveProcs) > 0 {
			processes = strings.Join(row.ActiveProcs, ", ")
		}
		events := "--"
		if len(row.Events) > 0 {
			events = strings.Join(row.Events, ", ")
		}
		tbl.Row(
			row.Time.Format("01/02 15:04"),
			fmt.Sprintf("%.0f%%", row.BatteryPct),
			power,
			temp,
			row.ChargeState,
			processes,
			events,
		)
	}

	return tbl.Render()
}
