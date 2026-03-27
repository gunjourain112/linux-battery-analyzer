package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

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
