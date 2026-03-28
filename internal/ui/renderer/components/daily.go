package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func Daily(records []domain.DailyRecord) string {
	if len(records) == 0 {
		return theme.Default.Subtle().Render("no daily records")
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers("Date", "Active", "Discharge", "Charge", "Avg W").
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
