package components

import (
	"strings"

	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
)

type processAggregate struct {
	Name    string
	Count   int
	CPUTime float64
	MemPeak float64
}

func aggregateProcesses(processes []domain.ProcessUsage) []processAggregate {
	groups := make(map[string]*processAggregate)
	for _, p := range processes {
		name := strings.TrimSpace(p.Name)
		if name == "" {
			name = strings.TrimSpace(p.RawName)
		}
		if name == "" {
			continue
		}
		g, ok := groups[name]
		if !ok {
			g = &processAggregate{Name: name}
			groups[name] = g
		}
		g.Count++
		g.CPUTime += p.CPUTime
		if p.MemPeak > g.MemPeak {
			g.MemPeak = p.MemPeak
		}
	}

	out := make([]processAggregate, 0, len(groups))
	for _, g := range groups {
		out = append(out, *g)
	}
	return out
}

func averageSessionDischargeRate(sessions []domain.Session) float64 {
	var totalHours float64
	var totalDrain float64
	for _, s := range sessions {
		rate := dischargeRate(s)
		if rate <= 0 {
			continue
		}
		totalHours += s.End.Sub(s.Start).Hours()
		totalDrain += s.StartPct - s.EndPct
	}
	if totalHours <= 0 {
		return 0
	}
	return totalDrain / totalHours
}

func topProcessImpact(impacts []domain.ProcessImpact) *domain.ProcessImpact {
	var top *domain.ProcessImpact
	for i := range impacts {
		imp := &impacts[i]
		if imp.DrainWatts <= 0 {
			continue
		}
		if top == nil || imp.DrainWatts > top.DrainWatts {
			top = imp
		}
	}
	return top
}

func dominantLoadBucket(profile domain.DischargeProfile) *domain.LoadBucket {
	var top *domain.LoadBucket
	for i := range profile.Buckets {
		b := &profile.Buckets[i]
		if b.Count == 0 {
			continue
		}
		if top == nil || b.Ratio > top.Ratio {
			top = b
		}
	}
	return top
}
