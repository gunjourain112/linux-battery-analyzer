package components

import (
	"math"
	"strings"
	"time"

	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
)

const graphWidth = 60
const graphHeight = 12

func BatteryGraph(points []domain.BatteryPoint) string {
	if len(points) < 2 {
		return "  (battery history data not available)"
	}

	sortPoints(points)
	start := points[0].Time
	end := points[len(points)-1].Time
	duration := end.Sub(start)
	if duration <= 0 {
		return "  (battery history data not available)"
	}

	colPct := make([]float64, graphWidth)
	for x := 0; x < graphWidth; x++ {
		frac := float64(x) / float64(graphWidth-1)
		targetTime := start.Add(time.Duration(float64(duration) * frac))
		colPct[x] = nearestBatteryPct(points, targetTime)
	}

	grid := make([][]rune, graphHeight)
	for i := range grid {
		grid[i] = make([]rune, graphWidth)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	for x := 0; x < graphWidth; x++ {
		pct := colPct[x]
		row := batteryPctToRow(pct)

		var ch rune
		if x < graphWidth-1 {
			diff := colPct[x+1] - pct
			if diff > 0.3 {
				ch = '╱'
			} else if diff < -0.3 {
				ch = '╲'
			} else {
				ch = '─'
			}
		} else {
			ch = '─'
		}

		if row >= 0 && row < graphHeight {
			grid[row][x] = ch
		}
	}

	labelRows := map[int]string{
		0:  "100%",
		2:  " 80%",
		4:  " 60%",
		6:  " 40%",
		9:  " 20%",
		11: "  0%",
	}

	var sb strings.Builder
	for r := 0; r < graphHeight; r++ {
		label := "    "
		if l, ok := labelRows[r]; ok {
			label = l
		}
		sb.WriteString(label + "│")
		sb.WriteString(string(grid[r]))
		sb.WriteString("\n")
	}
	sb.WriteString("    └" + strings.Repeat("─", graphWidth) + "\n")

	startLabel := start.Format("01/02 15:04")
	endLabel := end.Format("01/02 15:04")
	padding := graphWidth - len(startLabel) - len(endLabel)
	if padding < 1 {
		padding = 1
	}
	sb.WriteString("     " + startLabel + strings.Repeat(" ", padding) + endLabel)

	return sb.String()
}

func batteryPctToRow(pct float64) int {
	if pct > 100 {
		pct = 100
	}
	if pct < 0 {
		pct = 0
	}
	return int(math.Round((100 - pct) / 100 * float64(graphHeight-1)))
}

func nearestBatteryPct(points []domain.BatteryPoint, t time.Time) float64 {
	best := points[0].Percentage
	bestDelta := absGraphDur(points[0].Time.Sub(t))
	for _, p := range points[1:] {
		d := absGraphDur(p.Time.Sub(t))
		if d < bestDelta {
			bestDelta = d
			best = p.Percentage
		}
	}
	return best
}

func absGraphDur(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}

func sortPoints(points []domain.BatteryPoint) {
	for i := 1; i < len(points); i++ {
		j := i
		for j > 0 && points[j-1].Time.After(points[j].Time) {
			points[j-1], points[j] = points[j], points[j-1]
			j--
		}
	}
}
