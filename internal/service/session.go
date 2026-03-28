package service

import (
	"sort"
	"time"

	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
)

const sessionStartFuzz = time.Hour

func BuildSessions(events []domain.PowerEvent, points []domain.BatteryPoint, since, until time.Time) []domain.Session {
	sort.Slice(events, func(i, j int) bool { return events[i].Time.Before(events[j].Time) })

	var sessions []domain.Session
	var start time.Time

	if !since.IsZero() {
		start = since
	}

	for _, ev := range events {
		switch ev.Kind {
		case "sleep", "shutdown":
			if !start.IsZero() {
				s := domain.Session{Start: start, End: ev.Time}
				adjustSessionStart(&s, points)
				if s.End.Sub(s.Start) > time.Minute {
					sessions = append(sessions, s)
				}
				start = time.Time{}
			}
		case "resume", "boot":
			if start.IsZero() {
				start = ev.Time
			}
		}
	}

	if !start.IsZero() {
		end := time.Now()
		if !until.IsZero() && until.Before(end) {
			end = until
		}
		s := domain.Session{Start: start, End: end}
		adjustSessionStart(&s, points)
		sessions = append(sessions, s)
	}

	return sessions
}

func DischargeRate(s domain.Session) float64 {
	dur := s.End.Sub(s.Start).Hours()
	if dur <= 0 {
		return 0
	}
	drain := s.StartPct - s.EndPct
	if drain <= 0 {
		return 0
	}
	return drain / dur
}

func batteryAt(points []domain.BatteryPoint, t time.Time) float64 {
	best := 0.0
	bestDelta := time.Duration(1<<63 - 1)
	for _, p := range points {
		d := p.Time.Sub(t)
		if d < 0 {
			d = -d
		}
		if d < bestDelta {
			bestDelta = d
			best = p.Percentage
		}
	}
	return best
}

func adjustSessionStart(s *domain.Session, points []domain.BatteryPoint) {
	window := batteryWindow(points, s.Start, s.End)
	if len(window) > 0 {
		if window[0].Time.Sub(s.Start) > sessionStartFuzz {
			s.Start = window[0].Time
			if s.End.Sub(s.Start) <= time.Minute {
				return
			}
		}
		s.StartPct = window[0].Percentage
		s.EndPct = window[len(window)-1].Percentage
		return
	}

	s.StartPct = batteryAt(points, s.Start)
	s.EndPct = batteryAt(points, s.End)
}

func batteryWindow(points []domain.BatteryPoint, start, end time.Time) []domain.BatteryPoint {
	window := make([]domain.BatteryPoint, 0)
	for _, p := range points {
		if p.Time.Before(start) || p.Time.After(end) {
			continue
		}
		window = append(window, p)
	}
	sort.Slice(window, func(i, j int) bool { return window[i].Time.Before(window[j].Time) })
	return window
}
