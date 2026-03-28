package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/i18n"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func Thermals(tr i18n.Translator, stats domain.ThermalStats) string {
	if stats.Count == 0 {
		return theme.Default.Subtle().Render(tr.Get(i18n.NoThermalSamples))
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers(tr.Get(i18n.SamplesHeader), tr.Get(i18n.MinHeader), tr.Get(i18n.MaxHeader), tr.Get(i18n.AvgHeader)).
		StyleFunc(func(r, c int) lipgloss.Style {
			if r == -1 {
				return theme.Default.Header()
			}
			return theme.Default.Value().Align(lipgloss.Right)
		})

	tbl.Row(
		fmt.Sprintf("%d", stats.Count),
		fmt.Sprintf("%d C", stats.Min),
		fmt.Sprintf("%d C", stats.Max),
		fmt.Sprintf("%d C", stats.Avg),
	)

	return tbl.Render()
}
