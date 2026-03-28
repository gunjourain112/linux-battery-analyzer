package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/i18n"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func Daily(tr i18n.Translator, records []domain.DailyRecord) string {
	if len(records) == 0 {
		return theme.Default.Subtle().Render(tr.Get(i18n.NoDailyRecords))
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers(tr.Get(i18n.DateHeader), tr.Get(i18n.ActiveHeader), tr.Get(i18n.DrainHeader), tr.Get(i18n.ChargeHeader), tr.Get(i18n.AvgWHeader)).
		StyleFunc(func(r, c int) lipgloss.Style {
			if r == -1 {
				return theme.Default.Header()
			}
			if c == 0 {
				return theme.Default.Value()
			}
			return theme.Default.Value().Align(lipgloss.Right)
		})

	for _, r := range records {
		tbl.Row(
			r.Date.Format("01/02"),
			fmt.Sprintf("%dm", r.ActiveMin),
			fmt.Sprintf("%.0f%%", r.Discharge),
			fmt.Sprintf("%.0f%%", r.Charge),
			fmt.Sprintf("%.1f", r.AvgWatts),
		)
	}
	return tbl.Render()
}
