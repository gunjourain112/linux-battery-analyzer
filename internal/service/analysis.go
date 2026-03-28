package service

import (
	"sort"
	"time"

	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
)

const (
	profileThresholdIdle   = 3.0
	profileThresholdLight  = 7.0
	profileThresholdMedium = 12.0
	processMatchWindow     = 30 * time.Minute
)

func BuildDischargeProfile(ratePoints []domain.RatePoint) domain.DischargeProfile {
	buckets := []domain.LoadBucket{
		{Label: "idle (<3W)"},
		{Label: "light (3-7W)"},
		{Label: "medium (7-12W)"},
		{Label: "heavy (12W+)"},
	}

	if len(ratePoints) == 0 {
		return domain.DischargeProfile{Buckets: buckets}
	}

	points := make([]domain.RatePoint, 0, len(ratePoints))
	for _, p := range ratePoints {
		if p.State != "discharging" || p.Watts <= 0 {
			continue
		}
		points = append(points, p)
	}
	if len(points) == 0 {
		return domain.DischargeProfile{Buckets: buckets}
	}

	sort.Slice(points, func(i, j int) bool { return points[i].Time.Before(points[j].Time) })

	var totalWeight float64
	for i, p := range points {
		weight := dischargingSpan(points, i)
		if weight <= 0 {
			continue
		}

		idx := 0
		switch {
		case p.Watts < profileThresholdIdle:
			idx = 0
		case p.Watts < profileThresholdLight:
			idx = 1
		case p.Watts < profileThresholdMedium:
			idx = 2
		default:
			idx = 3
		}

		buckets[idx].Count++
		buckets[idx].AvgWatts += p.Watts * weight.Hours()
		buckets[idx].EstHours += weight.Hours()
		totalWeight += weight.Hours()
	}

	for i := range buckets {
		if buckets[i].EstHours > 0 {
			buckets[i].AvgWatts /= buckets[i].EstHours
			buckets[i].Ratio = buckets[i].EstHours / totalWeight * 100
		}
	}

	return domain.DischargeProfile{
		Buckets:    buckets,
		TotalCount: len(points),
	}
}

func dischargingSpan(points []domain.RatePoint, idx int) time.Duration {
	if len(points) == 1 {
		return 30 * time.Second
	}

	if idx < len(points)-1 {
		d := points[idx+1].Time.Sub(points[idx].Time)
		if d > 0 {
			return d
		}
	}

	if idx > 0 {
		d := points[idx].Time.Sub(points[idx-1].Time)
		if d > 0 {
			return d
		}
	}

	return 0
}

func BuildProcessImpacts(processes []domain.ProcessUsage, ratePoints []domain.RatePoint) []domain.ProcessImpact {
	impacts := make([]domain.ProcessImpact, 0, len(processes))
	for _, p := range processes {
		if p.Time.IsZero() {
			continue
		}
		watts := avgDischargingRate(ratePoints, p.Time, processMatchWindow)
		impacts = append(impacts, domain.ProcessImpact{
			Process:    p,
			DrainWatts: watts,
			Level:      classifyLoad(watts),
		})
	}

	sort.Slice(impacts, func(i, j int) bool {
		if impacts[i].DrainWatts == impacts[j].DrainWatts {
			return impacts[i].Process.MemPeak > impacts[j].Process.MemPeak
		}
		return impacts[i].DrainWatts > impacts[j].DrainWatts
	})
	return impacts
}

func avgDischargingRate(points []domain.RatePoint, t time.Time, window time.Duration) float64 {
	var sum float64
	var count int
	for _, p := range points {
		if p.State != "discharging" || p.Watts <= 0 {
			continue
		}
		d := p.Time.Sub(t)
		if d < 0 {
			d = -d
		}
		if d > window {
			continue
		}
		sum += p.Watts
		count++
	}
	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

func classifyLoad(watts float64) domain.LoadLevel {
	switch {
	case watts <= 0:
		return domain.LoadLevelUnknown
	case watts < profileThresholdIdle:
		return domain.LoadLevelLight
	case watts < profileThresholdLight:
		return domain.LoadLevelMedium
	default:
		return domain.LoadLevelHeavy
	}
}
