package service

import (
	"testing"
	"time"

	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
)

func TestBuildSessionsAdjustsPhantomBootStart(t *testing.T) {
	loc := time.Local
	boot := time.Date(2026, 3, 1, 8, 0, 0, 0, loc)
	sleep := time.Date(2026, 3, 1, 13, 0, 0, 0, loc)
	points := []domain.BatteryPoint{
		{Time: time.Date(2026, 3, 1, 10, 15, 0, 0, loc), Percentage: 92, State: "discharging"},
		{Time: time.Date(2026, 3, 1, 11, 0, 0, 0, loc), Percentage: 88, State: "discharging"},
		{Time: time.Date(2026, 3, 1, 12, 30, 0, 0, loc), Percentage: 81, State: "discharging"},
	}
	events := []domain.PowerEvent{
		{Time: boot, Kind: "boot"},
		{Time: sleep, Kind: "sleep"},
	}

	sessions := BuildSessions(events, points, time.Time{}, time.Time{})
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}

	got := sessions[0]
	wantStart := points[0].Time
	if !got.Start.Equal(wantStart) {
		t.Fatalf("expected start %s, got %s", wantStart, got.Start)
	}
	if !got.End.Equal(sleep) {
		t.Fatalf("expected end %s, got %s", sleep, got.End)
	}
	if got.StartPct != points[0].Percentage {
		t.Fatalf("expected start pct %.0f, got %.0f", points[0].Percentage, got.StartPct)
	}
	if got.EndPct != points[len(points)-1].Percentage {
		t.Fatalf("expected end pct %.0f, got %.0f", points[len(points)-1].Percentage, got.EndPct)
	}
}
