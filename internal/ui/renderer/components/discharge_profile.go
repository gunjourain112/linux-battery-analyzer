package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/i18n"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func DischargeProfile(tr i18n.Translator, profile domain.DischargeProfile) string {
	if profile.TotalCount == 0 {
		return theme.Default.Subtle().Render(tr.Get(i18n.NoDischargeProfile))
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers(tr.Get(i18n.BucketHeader), tr.Get(i18n.CountHeader), tr.Get(i18n.RatioHeader), tr.Get(i18n.AvgWHeader)).
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
