package components

import (
	"fmt"

	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/i18n"
)

func Summary(tr i18n.Translator, sessions []domain.Session, charging []domain.ChargingSession, events []domain.SystemEvent) string {
	var totalHours float64
	var totalDrain float64
	var worst domain.Session
	var worstRate float64
	var hasWorst bool

	for _, s := range sessions {
		rate := dischargeRate(s)
		if rate <= 0 {
			continue
		}
		totalHours += s.End.Sub(s.Start).Hours()
		totalDrain += s.StartPct - s.EndPct
		if !hasWorst || rate > worstRate {
			hasWorst = true
			worst = s
			worstRate = rate
		}
	}

	lines := []string{
		fmt.Sprintf("%s: %d", tr.Get(i18n.ReportSessions), len(sessions)),
		fmt.Sprintf("%s: %d", tr.Get(i18n.ReportCharging), len(charging)),
		fmt.Sprintf("%s: %d", tr.Get(i18n.ReportSystemEvents), len(events)),
	}
	if totalHours > 0 {
		lines = append(lines, fmt.Sprintf("%s: %.2f%%/h", tr.Get(i18n.AvgDischarge), totalDrain/totalHours))
	} else {
		lines = append(lines, fmt.Sprintf("%s: --", tr.Get(i18n.AvgDischarge)))
	}
	if hasWorst {
		lines = append(lines, fmt.Sprintf("%s: %s ~ %s (%.2f%%/h)",
			tr.Get(i18n.WorstSession),
			worst.Start.Format("01/02 15:04"),
			worst.End.Format("15:04"),
			worstRate,
		))
	}
	return joinLines(lines...)
}
