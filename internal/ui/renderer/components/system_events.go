package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func SystemEvents(events []domain.SystemEvent) string {
	if len(events) == 0 {
		return theme.Default.Subtle().Render("no system events")
	}

	limit := len(events)
	if limit > 8 {
		limit = 8
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers("Time", "Type", "Description").
		StyleFunc(func(r, c int) lipgloss.Style {
			if r == -1 {
				return theme.Default.Header()
			}
			return theme.Default.Value()
		})

	for i := 0; i < limit; i++ {
		ev := events[i]
		tbl.Row(
			ev.Time.Format("01/02 15:04"),
			ev.Type,
			ev.Desc,
		)
	}
	return tbl.Render()
}
