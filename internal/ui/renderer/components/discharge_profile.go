package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func DischargeProfile(profile domain.DischargeProfile) string {
	if profile.TotalCount == 0 {
		return theme.Default.Subtle().Render("no discharge profile")
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers("Bucket", "Count", "Ratio", "Avg W").
		StyleFunc(func(r, c int) lipgloss.Style {
			if r == -1 {
				return theme.Default.Header()
			}
			if c == 0 {
				return theme.Default.Label()
			}
			return theme.Default.Value().Align(lipgloss.Right)
		})

	for _, b := range profile.Buckets {
		if b.Count == 0 {
			continue
		}
		tbl.Row(
			b.Label,
			fmt.Sprintf("%d", b.Count),
			fmt.Sprintf("%.0f%%", b.Ratio),
			fmt.Sprintf("%.1f", b.AvgWatts),
		)
	}
	return tbl.Render()
}
