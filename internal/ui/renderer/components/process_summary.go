package components

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/i18n"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func ProcessSummary(tr i18n.Translator, processes []domain.ProcessUsage) string {
	aggs := aggregateProcesses(processes)
	if len(aggs) == 0 {
		return theme.Default.Subtle().Render(tr.Get(i18n.NoProcessSummaryData))
	}

	sort.Slice(aggs, func(i, j int) bool {
		if aggs[i].Count == aggs[j].Count {
			if aggs[i].CPUTime == aggs[j].CPUTime {
				return aggs[i].MemPeak > aggs[j].MemPeak
			}
			return aggs[i].CPUTime > aggs[j].CPUTime
		}
		return aggs[i].Count > aggs[j].Count
	})

	limit := len(aggs)
	if limit > 6 {
		limit = 6
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers(tr.Get(i18n.ProcessHeader), tr.Get(i18n.CountHeader), tr.Get(i18n.CPUHeader), tr.Get(i18n.MemHeader)).
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
		p := aggs[i]
		tbl.Row(
			p.Name,
			fmt.Sprintf("%d", p.Count),
			fmt.Sprintf("%.0fs", p.CPUTime),
			fmt.Sprintf("%.1f", p.MemPeak),
		)
	}

	return tbl.Render()
}
