package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type HardwareSpecs struct {
	OS      string
	Device  string
	CPU     string
	RAM     string
	Battery string
}

func (s HardwareSpecs) isEmpty() bool {
	return s.OS == "" && s.Device == "" && s.CPU == "" && s.RAM == "" && s.Battery == ""
}

var thermalLine = regexp.MustCompile(`\(([0-9]+) C\)`)

func loadThermalStats(since, until time.Time) (ThermalStats, error) {
	args := []string{"--no-pager", "-o", "short-iso"}
	if !since.IsZero() {
		args = append(args, "--since", since.Format("2006-01-02"))
	}
	if !until.IsZero() {
		args = append(args, "--until", until.Format("2006-01-02 23:59:59"))
	}

	out, err := exec.Command("journalctl", args...).Output()
	if err != nil {
		return ThermalStats{}, err
	}

	var stats ThermalStats
	sum := 0
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 25 || !strings.Contains(line, "Thermal Zone") {
			continue
		}
		match := thermalLine.FindStringSubmatch(line)
		if len(match) < 2 {
			continue
		}
		temp, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		if stats.Count == 0 || temp < stats.Min {
			stats.Min = temp
		}
		if stats.Count == 0 || temp > stats.Max {
			stats.Max = temp
		}
		stats.Count++
		sum += temp
	}
	if err := scanner.Err(); err != nil {
		return stats, err
	}
	if stats.Count > 0 {
		stats.Avg = sum / stats.Count
	}
	return stats, nil
}

func loadSpecs() HardwareSpecs {
	specs := HardwareSpecs{
		OS:      "--",
		Device:  "--",
		CPU:     "--",
		RAM:     "--",
		Battery: "--",
	}

	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				specs.OS = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), `"`)
				break
			}
		}
	}

	if data, err := os.ReadFile("/sys/class/dmi/id/product_name"); err == nil {
		if s := strings.TrimSpace(string(data)); s != "" {
			specs.Device = s
		}
	}

	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "model name") {
				if parts := strings.SplitN(line, ":", 2); len(parts) == 2 {
					specs.CPU = strings.TrimSpace(parts[1])
					break
				}
			}
		}
	}

	if data, err := os.ReadFile("/proc/meminfo"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "MemTotal:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					if kb, err := strconv.ParseFloat(fields[1], 64); err == nil {
						specs.RAM = fmt.Sprintf("%.1fGB", kb/1024/1024)
					}
				}
				break
			}
		}
	}

	batCmd := `upower -i $(upower -e | grep 'BAT' | head -1) 2>/dev/null | grep -E 'energy-full:|energy-full-design:|cycle-count:|capacity:'`
	if out, err := exec.Command("sh", "-c", batCmd).Output(); err == nil {
		var full, design, capPct, cycle string
		for _, line := range strings.Split(string(out), "\n") {
			fields := strings.Fields(strings.TrimSpace(line))
			if len(fields) < 2 {
				continue
			}
			switch {
			case strings.HasPrefix(line, "energy-full:"):
				full = fields[1]
			case strings.HasPrefix(line, "energy-full-design:"):
				design = fields[1]
			case strings.HasPrefix(line, "cycle-count:"):
				cycle = fields[1]
			case strings.HasPrefix(line, "capacity:"):
				capPct = strings.TrimSuffix(fields[1], "%")
			}
		}
		parts := []string{}
		if full != "" {
			parts = append(parts, "full "+full+"Wh")
		}
		if design != "" {
			parts = append(parts, "design "+design+"Wh")
		}
		if capPct != "" {
			parts = append(parts, capPct+"%")
		}
		if cycle != "" {
			parts = append(parts, "cycle "+cycle)
		}
		if len(parts) > 0 {
			specs.Battery = strings.Join(parts, ", ")
		}
	}

	return specs
}
