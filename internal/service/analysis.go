package service

import (
	"sort"
	"time"

	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
)

const (
	profileThresholdIdle   = 4.0
	profileThresholdLight  = 8.0
	profileThresholdMedium = 15.0
	processMatchWindow     = 30 * time.Minute
	timelineBucketSize     = 30 * time.Minute
)

func BuildDischargeProfile(ratePoints []domain.RatePoint) domain.DischargeProfile {
	buckets := []domain.LoadBucket{
		{Label: profileBucketLabel(0)},
		{Label: profileBucketLabel(1)},
		{Label: profileBucketLabel(2)},
		{Label: profileBucketLabel(3)},
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

func profileBucketLabel(idx int) string {
	switch idx {
	case 0:
		return "idle (<4W)"
	case 1:
		return "light (4-8W)"
	case 2:
		return "medium (8-15W)"
	default:
		return "heavy (15W+)"
	}
}

func BuildDetailedTimeline(points []domain.BatteryPoint, rates []domain.RatePoint, events []domain.SystemEvent, processes []domain.ProcessUsage, thermal domain.ThermalStats) []domain.DetailedTimelineRow {
	type bucket struct {
		t        time.Time
		lastPct  float64
		powerSum float64
		powerCnt int
		events   []string
		procs    []string
	}

	if len(points) == 0 {
		return nil
	}

	bucketMap := make(map[int64]*bucket)
	bucketKey := func(t time.Time) int64 {
		return t.Unix() / int64(timelineBucketSize.Seconds())
	}
	ensure := func(t time.Time) *bucket {
		key := bucketKey(t)
		if b, ok := bucketMap[key]; ok {
			return b
		}
		b := &bucket{t: t.Truncate(timelineBucketSize)}
		bucketMap[key] = b
		return b
	}

	sort.Slice(points, func(i, j int) bool { return points[i].Time.Before(points[j].Time) })
	for _, p := range points {
		b := ensure(p.Time)
		b.lastPct = p.Percentage
	}

	sort.Slice(rates, func(i, j int) bool { return rates[i].Time.Before(rates[j].Time) })
	for _, rp := range rates {
		if rp.State != "discharging" || rp.Watts <= 0 {
			continue
		}
		b := ensure(rp.Time)
		b.powerSum += rp.Watts
		b.powerCnt++
	}

	for _, ev := range events {
		b := ensure(ev.Time)
		switch ev.Type {
		case "sleep":
			b.events = append(b.events, "sleep")
		case "resume":
			b.events = append(b.events, "resume")
		case "shutdown":
			b.events = append(b.events, "shutdown")
		case "boot":
			b.events = append(b.events, "boot")
		}
	}

	for _, p := range processes {
		if p.Time.IsZero() {
			continue
		}
		b := ensure(p.Time)
		if len(b.procs) < 3 {
			b.procs = append(b.procs, p.Name)
		}
	}

	rows := make([]domain.DetailedTimelineRow, 0, len(bucketMap))
	for _, b := range bucketMap {
		row := domain.DetailedTimelineRow{
			Time:       b.t,
			BatteryPct: b.lastPct,
		}
		if b.powerCnt > 0 {
			row.PowerWatts = b.powerSum / float64(b.powerCnt)
		}
		row.ChargeState = dominantChargeState(points, b.t)
		row.PowerState = dominantPowerState(events, b.t)
		row.ActiveProcs = append(row.ActiveProcs, b.procs...)
		row.Events = append(row.Events, b.events...)
		if thermal.Count > 0 {
			// thermal summary is lightweight here; attach average later in renderer if needed
		}
		rows = append(rows, row)
	}

	sort.Slice(rows, func(i, j int) bool { return rows[i].Time.Before(rows[j].Time) })
	return rows
}

func dominantChargeState(points []domain.BatteryPoint, t time.Time) string {
	best := ""
	bestDelta := time.Duration(1<<63 - 1)
	for _, p := range points {
		d := absDur(p.Time.Sub(t))
		if d < bestDelta {
			bestDelta = d
			best = p.State
		}
	}
	return best
}

func dominantPowerState(events []domain.SystemEvent, t time.Time) string {
	best := ""
	bestDelta := time.Duration(1<<63 - 1)
	for _, ev := range events {
		d := absDur(ev.Time.Sub(t))
		if d < bestDelta {
			bestDelta = d
			best = ev.Type
		}
	}
	return best
}

func absDur(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
