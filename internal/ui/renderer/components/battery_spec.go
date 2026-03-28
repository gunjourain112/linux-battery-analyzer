package components

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
)

var batterySpecRe = regexp.MustCompile(`full\s+([0-9.]+)Wh(?:,\s*design\s+([0-9.]+)Wh)?(?:,\s*([0-9.]+)%)?`)

func parseBatterySpec(spec string) (fullWh float64, designWh float64, healthPct float64, ok bool) {
	m := batterySpecRe.FindStringSubmatch(spec)
	if len(m) < 2 {
		return 0, 0, 0, false
	}
	fullWh, _ = strconv.ParseFloat(m[1], 64)
	if len(m) > 2 && strings.TrimSpace(m[2]) != "" {
		designWh, _ = strconv.ParseFloat(m[2], 64)
	}
	if len(m) > 3 && strings.TrimSpace(m[3]) != "" {
		healthPct, _ = strconv.ParseFloat(m[3], 64)
	}
	if designWh <= 0 {
		designWh = fullWh
	}
	return fullWh, designWh, healthPct, fullWh > 0
}

func latestBatteryPoint(points []domain.BatteryPoint) (domain.BatteryPoint, bool) {
	if len(points) == 0 {
		return domain.BatteryPoint{}, false
	}
	latest := points[0]
	for _, p := range points[1:] {
		if p.Time.After(latest.Time) {
			latest = p
		}
	}
	return latest, true
}
