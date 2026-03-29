package service

import (
	"sort"
	"strings"
	"time"

	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
)

func BuildChargingSessions(points []domain.BatteryPoint, ratePoints []domain.RatePoint) []domain.ChargingSession {
	var sessions []domain.ChargingSession
	var cur *domain.ChargingSession

	for _, p := range points {
		if !isChargingState(p.State) {
			if cur != nil && cur.End.After(cur.Start) {
				applyChargingRates(cur, ratePoints)
				sessions = append(sessions, *cur)
			}
			cur = nil
			continue
		}

		if cur == nil {
			cur = &domain.ChargingSession{Start: p.Time, StartPct: p.Percentage}
		}
		cur.End = p.Time
		cur.EndPct = p.Percentage
	}

	if cur != nil && cur.End.After(cur.Start) {
		applyChargingRates(cur, ratePoints)
		sessions = append(sessions, *cur)
	}

	return sessions
}

func BuildDailySummary(sessions []domain.Session, charging []domain.ChargingSession) []domain.DailyRecord {
	type dayData struct {
		activeMin int
		discharge float64
		charge    float64
		sumWatts  float64
		samples   int
	}

	days := make(map[string]*dayData)
	touch := func(t time.Time) *dayData {
		key := t.Format("2006-01-02")
		if days[key] == nil {
			days[key] = &dayData{}
		}
		return days[key]
	}

	for _, s := range sessions {
		d := touch(s.Start)
		d.activeMin += int(s.End.Sub(s.Start).Minutes())
		drain := s.StartPct - s.EndPct
		if drain > 0 {
			d.discharge += drain
		}
		rate := DischargeRate(s)
		if rate > 0 {
			d.sumWatts += rate
			d.samples++
		}
	}

	for _, s := range charging {
		d := touch(s.Start)
		charge := s.EndPct - s.StartPct
		if charge > 0 {
			d.charge += charge
		}
		if s.AvgChargeW > 0 {
			d.sumWatts += s.AvgChargeW
			d.samples++
		}
	}

	keys := make([]string, 0, len(days))
	for k := range days {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make([]domain.DailyRecord, 0, len(keys))
	for _, k := range keys {
		d, _ := days[k]
		rec := domain.DailyRecord{}
		rec.Date, _ = time.ParseInLocation("2006-01-02", k, time.Local)
		rec.ActiveMin = d.activeMin
		rec.Discharge = d.discharge
		rec.Charge = d.charge
		if d.samples > 0 {
			rec.AvgWatts = d.sumWatts / float64(d.samples)
		}
		out = append(out, rec)
	}
	return out
}

func BuildSystemEvents(powerEvents []domain.PowerEvent, points []domain.BatteryPoint) []domain.SystemEvent {
	var events []domain.SystemEvent

	for _, ev := range powerEvents {
		switch ev.Kind {
		case "sleep":
			events = append(events, domain.SystemEvent{Time: ev.Time, Type: "sleep", Desc: "sleep"})
		case "resume":
			events = append(events, domain.SystemEvent{Time: ev.Time, Type: "resume", Desc: "resume"})
		case "shutdown":
			events = append(events, domain.SystemEvent{Time: ev.Time, Type: "shutdown", Desc: "shutdown"})
		case "boot":
			events = append(events, domain.SystemEvent{Time: ev.Time, Type: "boot", Desc: "boot"})
		}
	}

	var prevState string
	for _, p := range points {
		cur := batteryStateGroup(p.State)
		if cur == "" {
			continue
		}
		if prevState != "" && cur != prevState {
			switch {
			case prevState == "charging" && cur == "discharging":
				events = append(events, domain.SystemEvent{Time: p.Time, Type: "ac-unplugged", Desc: "ac unplugged"})
			case prevState == "discharging" && cur == "charging":
				events = append(events, domain.SystemEvent{Time: p.Time, Type: "ac-plugged", Desc: "ac plugged in"})
			}
		}
		prevState = cur
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Time.Before(events[j].Time)
	})

	if len(events) == 0 {
		return events
	}

	deduped := make([]domain.SystemEvent, 0, len(events))
	for _, ev := range events {
		if len(deduped) > 0 {
			last := deduped[len(deduped)-1]
			if last.Type == ev.Type && last.Time.Truncate(time.Minute).Equal(ev.Time.Truncate(time.Minute)) {
				continue
			}
		}
		deduped = append(deduped, ev)
	}
	return deduped
}

func applyChargingRates(sess *domain.ChargingSession, ratePoints []domain.RatePoint) {
	var sum float64
	var count int
	for _, p := range ratePoints {
		if p.State != "charging" || p.Watts <= 0 {
			continue
		}
		if p.Time.Before(sess.Start) || p.Time.After(sess.End) {
			continue
		}
		sum += p.Watts
		count++
		if p.Watts > sess.PeakChargeW {
			sess.PeakChargeW = p.Watts
		}
	}
	if count > 0 {
		sess.AvgChargeW = sum / float64(count)
	}
}

func isChargingState(state string) bool {
	switch strings.ToLower(strings.TrimSpace(state)) {
	case "charging", "fully-charged":
		return true
	default:
		return false
	}
}

func batteryStateGroup(state string) string {
	switch strings.ToLower(strings.TrimSpace(state)) {
	case "charging", "fully-charged":
		return "charging"
	case "discharging":
		return "discharging"
	default:
		return ""
	}
}
