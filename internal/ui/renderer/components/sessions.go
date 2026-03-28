package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func Sessions(sessions []domain.Session) string {
	if len(sessions) == 0 {
		return theme.Default.Subtle().Render("no sessions")
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers("Start", "End", "Duration", "Drain", "Rate").
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
		rate := dischargeRate(s)
		drain := s.StartPct - s.EndPct
		tbl.Row(
			s.Start.Format("01/02 15:04"),
			s.End.Format("01/02 15:04"),
			formatDuration(s.End.Sub(s.Start)),
			fmt.Sprintf("%.0f%%", drain),
			fmt.Sprintf("%.2f%%/h", rate),
		)
	}
	return tbl.Render()
}
