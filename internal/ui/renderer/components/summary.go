package components

import (
	"fmt"

	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
)

func Summary(sessions []domain.Session, charging []domain.ChargingSession, events []domain.SystemEvent) string {
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
		fmt.Sprintf("sessions: %d", len(sessions)),
		fmt.Sprintf("charging: %d", len(charging)),
		fmt.Sprintf("events: %d", len(events)),
	}
	if totalHours > 0 {
		lines = append(lines, fmt.Sprintf("avg discharge: %.2f%%/h", totalDrain/totalHours))
	} else {
		lines = append(lines, "avg discharge: --")
	}
	if hasWorst {
		lines = append(lines, fmt.Sprintf("worst session: %s ~ %s (%.2f%%/h)",
			worst.Start.Format("01/02 15:04"),
			worst.End.Format("15:04"),
			worstRate,
		))
	}
	return joinLines(lines...)
}
