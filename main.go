package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

	fmt.Println("=== battery ===")
	for _, p := range points {
		fmt.Printf("%s  %5.1f%%  %s\n", p.Time.Format("2006-01-02 15:04"), p.Percentage, p.State)
	}

	fmt.Println("\n=== power events ===")
	for _, e := range events {
		fmt.Printf("%s  %s\n", e.Time.Format("2006-01-02 15:04"), e.Kind)
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
