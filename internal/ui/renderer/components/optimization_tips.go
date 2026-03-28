package components

import (
	"fmt"
	"strings"

	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/i18n"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func OptimizationTips(tr i18n.Translator, sessions []domain.Session, impacts []domain.ProcessImpact, profile domain.DischargeProfile, thermal domain.ThermalStats) string {
	tips := make([]string, 0, 4)

	if top := topProcessImpact(impacts); top != nil && top.DrainWatts > 0 {
		tips = append(tips, fmt.Sprintf(tr.Get(i18n.TipTopProcess), top.Process.Name, top.DrainWatts))
	}

	if bucket := dominantLoadBucket(profile); bucket != nil && bucket.Ratio >= 35 && strings.Contains(strings.ToLower(bucket.Label), "heavy") {
		tips = append(tips, fmt.Sprintf(tr.Get(i18n.TipHeavyLoad), bucket.Ratio))
	}

	if thermal.Count > 0 && thermal.Max > 0 {
		tips = append(tips, fmt.Sprintf(tr.Get(i18n.TipThermal), thermal.Max))
	}

	if avg := averageSessionDischargeRate(sessions); avg > 0 {
		tips = append(tips, fmt.Sprintf(tr.Get(i18n.TipDrainRate), avg))
	}

	if len(tips) == 0 {
		return theme.Default.Subtle().Render(tr.Get(i18n.TipNoObviousIssue))
	}

	var b strings.Builder
	for _, tip := range tips {
		b.WriteString("  • ")
		b.WriteString(theme.Default.Value().Render(tip))
		b.WriteString("\n")
	}
	return strings.TrimRight(b.String(), "\n")
}
