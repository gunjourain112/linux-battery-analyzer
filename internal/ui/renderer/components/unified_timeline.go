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

func UnifiedTimeline(rows []domain.DetailedTimelineRow, tr *i18n.Translator) string {
	t := theme.Default
	if len(rows) == 0 {
		return t.Subtle().Render(tr.Get(i18n.NoTimelineData))
	}

	tt := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(t.Subtle()).
		Headers(
			tr.Get(i18n.TimeHeader),
			tr.Get(i18n.BatteryHeader),
			tr.Get(i18n.HeatHeader),
			tr.Get(i18n.ChargeHeader),
			tr.Get(i18n.PowerHeader),
			tr.Get(i18n.EventHeader),
		).
		StyleFunc(func(r, c int) lipgloss.Style {
			if r == -1 {
				return t.Header()
			}
			return t.Value()
		})

	for _, row := range rows {
		powerStr := "--"
		if row.PowerWatts > 0 {
			powerStr = fmt.Sprintf("%.1fW", row.PowerWatts)
		}

		tt.Row(
			row.Time.Format("01/02 15:04"),
			fmt.Sprintf("%.0f%%", row.BatteryPct),
			"--",
			chargeStateCompact(row.ChargeState),
			powerStr,
			strings.Join(row.Events, ", "),
		)
	}
	return tt.Render()
}

func chargeStateCompact(state string) string {
	switch state {
	case "charging":
		return "⚡ 충전"
	case "discharging":
		return "🔋 방전"
	case "fully-charged":
		return "● 완충"
	default:
		return "--"
	}
}
