package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func Thermals(stats domain.ThermalStats) string {
	if stats.Count == 0 {
		return theme.Default.Subtle().Render("no thermal samples")
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers("Samples", "Min", "Max", "Avg").
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
