package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type BatteryPoint struct {
	Time       time.Time
	Percentage float64
	State      string
}

type PowerEvent struct {
	Time time.Time
	Kind string // sleep, resume, shutdown, boot
}

type Session struct {
	Start    time.Time
	End      time.Time
	StartPct float64
	EndPct   float64
}

func main() {
	var since, until time.Time

	if len(os.Args) >= 3 {
		var err error
		since, err = time.ParseInLocation("2006-01-02", os.Args[1], time.Local)
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid since date:", os.Args[1])
			os.Exit(1)
		}
		until, err = time.ParseInLocation("2006-01-02", os.Args[2], time.Local)
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid until date:", os.Args[2])
			os.Exit(1)
		}
		until = until.Add(24*time.Hour - time.Second)
	}

	points, err := loadHistory(since, until)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load history:", err)
		os.Exit(1)
	}

	events, err := loadPowerEvents(since, until)
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning: could not load journal:", err)
	}

	sessions := buildSessions(events, points, since, until)

	fmt.Println("=== sessions ===")
	for i, s := range sessions {
		dur := s.End.Sub(s.Start)
		drain := s.StartPct - s.EndPct
		rate := dischargeRate(s)
		fmt.Printf("[%d] %s ~ %s  (%dh%02dm)  %.0f%% → %.0f%%  (%.1f%% drain)\n",
			i+1,
			s.Start.Format("01/02 15:04"),
			s.End.Format("15:04"),
			int(dur.Hours()), int(dur.Minutes())%60,
			s.StartPct, s.EndPct, drain,
		)
		if rate > 0 {
			fmt.Printf("    rate: %.2f%%/h\n", rate)
		}
	}

	printSummary(sessions)
}

func buildSessions(events []PowerEvent, points []BatteryPoint, since, until time.Time) []Session {
	// 이벤트 + 배터리 포인트 시간 기준 정렬
	sort.Slice(events, func(i, j int) bool { return events[i].Time.Before(events[j].Time) })

	var sessions []Session
	var start time.Time

	if !since.IsZero() {
		start = since
	}

	for _, ev := range events {
		switch ev.Kind {
		case "sleep", "shutdown":
			if !start.IsZero() {
				s := Session{Start: start, End: ev.Time}
				s.StartPct = batteryAt(points, start)
				s.EndPct = batteryAt(points, ev.Time)
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

	// 마지막 세션이 until까지 열려있으면 닫기
	if !start.IsZero() {
		end := time.Now()
		if !until.IsZero() && until.Before(end) {
			end = until
		}
		s := Session{Start: start, End: end}
		s.StartPct = batteryAt(points, start)
		s.EndPct = batteryAt(points, end)
		sessions = append(sessions, s)
	}

	return sessions
}

func batteryAt(points []BatteryPoint, t time.Time) float64 {
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

func dischargeRate(s Session) float64 {
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

func printSummary(sessions []Session) {
	var totalHours float64
	var totalDrain float64
	var worst Session
	var worstRate float64
	var hasWorst bool

	for _, s := range sessions {
		rate := dischargeRate(s)
		if rate <= 0 {
			continue
		}

		hours := s.End.Sub(s.Start).Hours()
		totalHours += hours
		totalDrain += s.StartPct - s.EndPct

		if !hasWorst || rate > worstRate {
			hasWorst = true
			worst = s
			worstRate = rate
		}
	}

	fmt.Println()
	fmt.Println("=== summary ===")
	fmt.Printf("sessions: %d\n", len(sessions))

	if totalHours == 0 {
		fmt.Println("avg discharge: --")
		fmt.Println("worst session: --")
		return
	}

	avgRate := totalDrain / totalHours
	fmt.Printf("avg discharge: %.2f%%/h\n", avgRate)

	if hasWorst {
		dur := worst.End.Sub(worst.Start)
		fmt.Printf(
			"worst session: %s ~ %s  (%dh%02dm)  %.2f%%/h\n",
			worst.Start.Format("01/02 15:04"),
			worst.End.Format("15:04"),
			int(dur.Hours()), int(dur.Minutes())%60,
			worstRate,
		)
	} else {
		fmt.Println("worst session: --")
	}
}

func loadPowerEvents(since, until time.Time) ([]PowerEvent, error) {
	args := []string{"--no-pager", "-o", "short-iso"}
	if !since.IsZero() {
		args = append(args, "--since", since.Format("2006-01-02"))
	}
	if !until.IsZero() {
		args = append(args, "--until", until.Format("2006-01-02 23:59:59"))
	}

	out, err := exec.Command("journalctl", args...).Output()
	if err != nil {
		return nil, err
	}

	var events []PowerEvent
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 25 {
			continue
		}
		t, err := time.ParseInLocation("2006-01-02T15:04:05-0700", line[:25], time.Local)
		if err != nil {
			continue
		}
		low := strings.ToLower(line)
		var kind string
		switch {
		case strings.Contains(low, "suspending system"),
			strings.Contains(low, "pm: suspend entry"):
			kind = "sleep"
		case strings.Contains(low, "system resumed"),
			strings.Contains(low, "pm: early resume"):
			kind = "resume"
		case strings.Contains(low, "powering off"),
			strings.Contains(low, "reached target power-off"):
			kind = "shutdown"
		case strings.Contains(low, "startup finished in"):
			kind = "boot"
		}
		if kind == "" {
			continue
		}
		events = append(events, PowerEvent{Time: t, Kind: kind})
	}
	return events, nil
}

func loadHistory(since, until time.Time) ([]BatteryPoint, error) {
	files, err := filepath.Glob("/var/lib/upower/history-charge-*.dat")
	if err != nil || len(files) == 0 {
		return nil, fmt.Errorf("no upower history files found")
	}

	var points []BatteryPoint
	for _, f := range files {
		pts, err := parseHistoryFile(f, since, until)
		if err != nil {
			continue
		}
		points = append(points, pts...)
	}
	return points, nil
}

func parseHistoryFile(path string, since, until time.Time) ([]BatteryPoint, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var points []BatteryPoint
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) < 3 {
			continue
		}
		epoch, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			continue
		}
		pct, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			continue
		}
		if parts[2] == "unknown" || pct <= 0 {
			continue
		}
		t := time.Unix(epoch, 0)
		if !since.IsZero() && t.Before(since) {
			continue
		}
		if !until.IsZero() && t.After(until) {
			continue
		}
		points = append(points, BatteryPoint{
			Time:       t,
			Percentage: pct,
			State:      parts[2],
		})
	}
	return points, scanner.Err()
}
