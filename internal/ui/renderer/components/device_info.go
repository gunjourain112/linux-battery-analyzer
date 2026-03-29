package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/i18n"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func DeviceInfo(tr i18n.Translator, config domain.Config, specs domain.HardwareSpecs, sessions []domain.Session, battery []domain.BatteryPoint, profile domain.DischargeProfile, thermal domain.ThermalStats) string {
	sections := []struct {
		title string
		body  string
	}{
		{title: "1-A. " + tr.Get(i18n.DeviceSpecsSection), body: Specs(tr, specs)},
		{title: "1-B. " + tr.Get(i18n.BatteryHealthSection), body: BatteryHealth(specs, battery, tr)},
		{title: "1-C. " + tr.Get(i18n.AnalysisSummarySection), body: AnalysisSummary(tr, config, specs, sessions, battery, profile, thermal)},
	}

	var b strings.Builder
	for i, sec := range sections {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(theme.Default.Label().Render(sec.title))
		b.WriteString("\n\n")
		b.WriteString(sec.body)
		if !strings.HasSuffix(sec.body, "\n") {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func AnalysisSummary(tr i18n.Translator, config domain.Config, specs domain.HardwareSpecs, sessions []domain.Session, battery []domain.BatteryPoint, profile domain.DischargeProfile, thermal domain.ThermalStats) string {
	current, ok := latestBatteryPoint(battery)
	if !ok {
		return theme.Default.Subtle().Render(tr.Get(i18n.NoBatteryHealthData))
	}

	rows := [][]string{
		{tr.Get(i18n.AnalysisPeriodHeader), fmt.Sprintf("%s ~ %s", config.Since.Format("2006-01-02 15:04:05"), config.Until.Format("2006-01-02 15:04:05"))},
		{tr.Get(i18n.ActualUseHeader), formatDurationCard(totalSessionDuration(sessions))},
		{tr.Get(i18n.BatteryStateHeader), fmt.Sprintf("%.0f%% %s", current.Percentage, batteryStateLabel(tr, current.State))},
	}
	if avg := averageSessionDischargeRate(sessions); avg > 0 {
		rows = append(rows, []string{tr.Get(i18n.AvgLoadHeader), fmt.Sprintf("%.2f%%/h", avg)})
	}
	if thermal.Count > 0 {
		rows = append(rows, []string{tr.Get(i18n.TempRangeHeader), fmt.Sprintf("%d / %d / %d", thermal.Min, thermal.Max, thermal.Avg)})
	}

	if estimate := analysisRemainingEstimate(specs, current, profile); estimate != "" {
		rows = append(rows, []string{tr.Get(i18n.ExpectedRemainHeader), estimate})
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		StyleFunc(func(r, c int) lipgloss.Style {
			if r == -1 {
				return theme.Default.Header()
			}
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

func analysisRemainingEstimate(specs domain.HardwareSpecs, current domain.BatteryPoint, profile domain.DischargeProfile) string {
	_, designWh, _, ok := parseBatterySpec(specs.Battery)
	if !ok || designWh <= 0 {
		return ""
	}

	watts := weightedAverageWatts(profile)
	if watts <= 0 {
		return ""
	}

	remainingWh := designWh * current.Percentage / 100
	return fmt.Sprintf("약 %s", formatDurationFromHours(remainingWh/watts))
}
