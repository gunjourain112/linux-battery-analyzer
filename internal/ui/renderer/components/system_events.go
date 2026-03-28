package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/i18n"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func SystemEvents(tr i18n.Translator, events []domain.SystemEvent) string {
	if len(events) == 0 {
		return theme.Default.Subtle().Render(tr.Get(i18n.NoSystemEvents))
	}

	limit := len(events)
	if limit > 8 {
		limit = 8
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers(tr.Get(i18n.TimeHeader), tr.Get(i18n.TypeHeader), tr.Get(i18n.DescriptionHeader)).
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
