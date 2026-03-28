package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/i18n"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func InsightDashboard(tr i18n.Translator, sessions []domain.Session, processes []domain.ProcessUsage, impacts []domain.ProcessImpact, profile domain.DischargeProfile, thermal domain.ThermalStats) string {
	rows := make([][]string, 0, 6)

	if len(sessions) > 0 {
		rows = append(rows, []string{tr.Get(i18n.ReportSessions), fmt.Sprintf("%d", len(sessions))})
	}
	if len(processes) > 0 {
		rows = append(rows, []string{tr.Get(i18n.ProcessSamplesLabel), fmt.Sprintf("%d", len(processes))})
	}
	if top := topProcessImpact(impacts); top != nil && top.DrainWatts > 0 {
		rows = append(rows, []string{tr.Get(i18n.TopDrainLabel), fmt.Sprintf("%s (%.1fW)", top.Process.Name, top.DrainWatts)})
	}
	if bucket := dominantLoadBucket(profile); bucket != nil {
		rows = append(rows, []string{tr.Get(i18n.HeavyLoadLabel), fmt.Sprintf("%s %.0f%%", bucket.Label, bucket.Ratio)})
	}
	if thermal.Count > 0 && thermal.Max > 0 {
		rows = append(rows, []string{tr.Get(i18n.PeakTempLabel), fmt.Sprintf("%d°C", thermal.Max)})
	}
	if avg := averageSessionDischargeRate(sessions); avg > 0 {
		rows = append(rows, []string{tr.Get(i18n.AvgDischarge), fmt.Sprintf("%.2f%%/h", avg)})
	}

	if len(rows) == 0 {
		return theme.Default.Subtle().Render(tr.Get(i18n.NoInsightDashboardData))
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		StyleFunc(func(r, c int) lipgloss.Style {
			if c == 0 {
				return theme.Default.Label()
			}
			return theme.Default.Value()
		})

	for _, row := range rows {
		tbl.Row(row...)
	}
	return tbl.Render()
}
