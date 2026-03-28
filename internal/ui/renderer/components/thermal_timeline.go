package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/i18n"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func ThermalTimeline(tr i18n.Translator, snapshots []domain.ThermalSnapshot) string {
	if len(snapshots) == 0 {
		return theme.Default.Subtle().Render(tr.Get(i18n.NoThermalSamples))
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers(tr.Get(i18n.TimeHeader), tr.Get(i18n.MinHeader), tr.Get(i18n.AvgHeader), tr.Get(i18n.MaxHeader), tr.Get(i18n.SamplesHeader)).
		StyleFunc(func(r, c int) lipgloss.Style {
			if r == -1 {
				return theme.Default.Header()
			}
			return theme.Default.Value().Align(lipgloss.Right)
		})

	limit := len(snapshots)
	if limit > 12 {
		limit = 12
	}
	for i := 0; i < limit; i++ {
		s := snapshots[i]
		tbl.Row(
			s.Hour,
			fmt.Sprintf("%d°C", s.Min),
			fmt.Sprintf("%d°C", s.Avg),
			fmt.Sprintf("%d°C", s.Max),
			fmt.Sprintf("%d", s.Count),
		)
	}
	return tbl.Render()
}
