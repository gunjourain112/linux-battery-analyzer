package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func ProcessImpacts(impacts []domain.ProcessImpact) string {
	if len(impacts) == 0 {
		return theme.Default.Subtle().Render("no process impact data")
	}

	limit := len(impacts)
	if limit > 5 {
		limit = 5
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers("Process", "Drain W", "Level", "CPU s", "Mem M").
		StyleFunc(func(r, c int) lipgloss.Style {
			if r == -1 {
				return theme.Default.Header()
			}
			if c == 0 {
				return theme.Default.Value()
			}
			return theme.Default.Value().Align(lipgloss.Right)
		})

	for i := 0; i < limit; i++ {
		p := impacts[i]
		tbl.Row(
			p.Process.Name,
			fmt.Sprintf("%.1f", p.DrainWatts),
			levelString(p.Level),
			fmt.Sprintf("%.0f", p.Process.CPUTime),
			fmt.Sprintf("%.1f", p.Process.MemPeak),
		)
	}
	return tbl.Render()
}
