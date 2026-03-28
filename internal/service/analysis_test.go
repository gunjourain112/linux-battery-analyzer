package service

import (
	"testing"
	"time"

	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
)

func TestBuildDischargeProfileUsesTimeWeighting(t *testing.T) {
	loc := time.Local
	points := []domain.RatePoint{
		{Time: time.Date(2026, 3, 1, 9, 0, 0, 0, loc), Watts: 2.0, State: "discharging"},
		{Time: time.Date(2026, 3, 1, 9, 25, 0, 0, loc), Watts: 2.5, State: "discharging"},
		{Time: time.Date(2026, 3, 1, 9, 50, 0, 0, loc), Watts: 2.8, State: "discharging"},
		{Time: time.Date(2026, 3, 1, 10, 0, 0, 0, loc), Watts: 14.0, State: "discharging"},
		{Time: time.Date(2026, 3, 1, 10, 2, 0, 0, loc), Watts: 15.0, State: "discharging"},
		{Time: time.Date(2026, 3, 1, 10, 4, 0, 0, loc), Watts: 16.0, State: "discharging"},
		{Time: time.Date(2026, 3, 1, 10, 6, 0, 0, loc), Watts: 15.5, State: "discharging"},
		{Time: time.Date(2026, 3, 1, 10, 8, 0, 0, loc), Watts: 14.5, State: "discharging"},
		{Time: time.Date(2026, 3, 1, 10, 10, 0, 0, loc), Watts: 15.0, State: "discharging"},
	}

	profile := BuildDischargeProfile(points)
	if profile.TotalCount != len(points) {
		t.Fatalf("expected total count %d, got %d", len(points), profile.TotalCount)
	}

	light := profile.Buckets[0]
	heavy := profile.Buckets[3]
	if light.Ratio <= heavy.Ratio {
		t.Fatalf("expected light bucket ratio > heavy bucket ratio, got light=%.2f heavy=%.2f", light.Ratio, heavy.Ratio)
	}
	if light.AvgWatts >= heavy.AvgWatts {
		t.Fatalf("expected heavy bucket avg watts to be higher, got light=%.2f heavy=%.2f", light.AvgWatts, heavy.AvgWatts)
	}
	if light.EstHours <= heavy.EstHours {
		t.Fatalf("expected light bucket duration to be higher, got light=%.2f heavy=%.2f", light.EstHours, heavy.EstHours)
	}
}
