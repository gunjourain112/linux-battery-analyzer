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

func main() {
	points, err := loadHistory()
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

func loadHistory() ([]BatteryPoint, error) {
	files, err := filepath.Glob("/var/lib/upower/history-charge-*.dat")
	if err != nil || len(files) == 0 {
		return nil, fmt.Errorf("no upower history files found")
	}

	var points []BatteryPoint
	for _, f := range files {
		pts, err := parseHistoryFile(f)
		if err != nil {
			continue
		}
		points = append(points, pts...)
	}
	return points, nil
}

func parseHistoryFile(path string) ([]BatteryPoint, error) {
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
		points = append(points, BatteryPoint{
			Time:       time.Unix(epoch, 0),
			Percentage: pct,
			State:      parts[2],
		})
	}
	return points, scanner.Err()
}

func getBatteryDevice() (string, error) {
	out, err := exec.Command("upower", "-e").Output()
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if strings.Contains(line, "BAT") {
			return line, nil
		}
	}
	return "", fmt.Errorf("no battery found")
}
