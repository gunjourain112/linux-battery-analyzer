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
		if aggs[i].MemPeak == aggs[j].MemPeak {
			if aggs[i].CPUTime == aggs[j].CPUTime {
				return aggs[i].Count > aggs[j].Count
			}
			return aggs[i].CPUTime > aggs[j].CPUTime
		}
		return aggs[i].MemPeak > aggs[j].MemPeak
	})

	limit := len(aggs)
	if limit > 20 {
		limit = 20
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
			formatCPUTime(p.CPUTime),
			formatMemoryPeak(p.MemPeak),
		)
	}

	return tbl.Render()
}

func formatCPUTime(seconds float64) string {
	if seconds <= 0 {
		return "--"
	}
	h := int(seconds) / 3600
	m := (int(seconds) % 3600) / 60
	s := int(seconds) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %02dm %02ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %02ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

func formatMemoryPeak(mb float64) string {
	if mb <= 0 {
		return "--"
	}
	if mb >= 1024 {
		return fmt.Sprintf("%.1fG", mb/1024)
	}
	return fmt.Sprintf("%.1fM", mb)
}
