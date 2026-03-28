package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/i18n"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func Charging(tr i18n.Translator, sessions []domain.ChargingSession) string {
	if len(sessions) == 0 {
		return theme.Default.Subtle().Render(tr.Get(i18n.NoChargingSessions))
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers(tr.Get(i18n.StartHeader), tr.Get(i18n.EndHeader), tr.Get(i18n.DurationHeader), tr.Get(i18n.ChargeHeader), tr.Get(i18n.AvgWHeader), tr.Get(i18n.PeakWHeader)).
		StyleFunc(func(r, c int) lipgloss.Style {
			if r == -1 {
				return theme.Default.Header()
			}
			if c >= 2 {
				return theme.Default.Value().Align(lipgloss.Right)
			}
			return theme.Default.Value()
		})

	for _, s := range sessions {
		gain := s.EndPct - s.StartPct
		tbl.Row(
			s.Start.Format("01/02 15:04"),
			s.End.Format("01/02 15:04"),
			formatDuration(s.End.Sub(s.Start)),
			fmt.Sprintf("%.0f%%", gain),
			fmt.Sprintf("%.1f", s.AvgChargeW),
			fmt.Sprintf("%.1f", s.PeakChargeW),
		)
	}
	return tbl.Render()
}
