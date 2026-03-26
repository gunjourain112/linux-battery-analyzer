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

type BatteryPoint struct {
	Time       time.Time
	Percentage float64
	State      string
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

	if len(points) == 0 {
		fmt.Println("no battery history found")
		return
	}

	for _, p := range points {
		fmt.Printf("%s  %5.1f%%  %s\n", p.Time.Format("2006-01-02 15:04"), p.Percentage, p.State)
	}
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
