package infrastructure

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
)

var (
	resourceLine = regexp.MustCompile(`^.*?: (.+): Consumed (.*) CPU time, (.*) memory peak\.$`)
	thermalLine  = regexp.MustCompile(`\(([0-9]+) C\)`)
)

func LoadHistory(since, until time.Time) ([]domain.BatteryPoint, error) {
	files, err := filepathGlob("/var/lib/upower/history-charge-*.dat")
	if err != nil || len(files) == 0 {
		return nil, fmt.Errorf("no upower history files found")
	}

	var points []domain.BatteryPoint
	for _, f := range files {
		pts, err := parseHistoryFile(f, since, until)
		if err != nil {
			continue
		}
		points = append(points, pts...)
	}
	return points, nil
}

func LoadRateHistory(since, until time.Time) ([]domain.RatePoint, error) {
	files, err := filepathGlob("/var/lib/upower/history-rate-*.dat")
	if err != nil || len(files) == 0 {
		return nil, nil
	}

	var points []domain.RatePoint
	for _, f := range files {
		pts, err := parseRateFile(f, since, until)
		if err != nil {
			continue
		}
		points = append(points, pts...)
	}
	return points, nil
}

func LoadPowerEvents(since, until time.Time) ([]domain.PowerEvent, error) {
	args := []string{"--no-pager", "-o", "short-iso"}
	if !since.IsZero() {
		args = append(args, "--since", since.Format("2006-01-02"))
	}
	if !until.IsZero() {
		args = append(args, "--until", until.Format("2006-01-02 23:59:59"))
	}

	out, err := exec.Command("journalctl", args...).CombinedOutput()
	if err != nil && len(strings.TrimSpace(string(out))) == 0 {
		return nil, err
	}

	var events []domain.PowerEvent
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
		events = append(events, domain.PowerEvent{Time: t, Kind: kind})
	}
	return events, nil
}

func LoadProcessUsage(since, until time.Time) ([]domain.ProcessUsage, error) {
	args := []string{"--no-pager", "-o", "short-iso"}
	if !since.IsZero() {
		args = append(args, "--since", since.Format("2006-01-02"))
	}
	if !until.IsZero() {
		args = append(args, "--until", until.Format("2006-01-02 23:59:59"))
	}

	out, err := exec.Command("journalctl", args...).CombinedOutput()
	if err != nil && len(strings.TrimSpace(string(out))) == 0 {
		return nil, err
	}

	var processes []domain.ProcessUsage
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 25 || !strings.Contains(line, "memory peak") {
			continue
		}

		match := resourceLine.FindStringSubmatch(line)
		if len(match) < 4 {
			continue
		}

		t, err := time.ParseInLocation("2006-01-02T15:04:05-0700", line[:25], time.Local)
		if err != nil {
			continue
		}

		processes = append(processes, domain.ProcessUsage{
			Time:    t,
			Name:    normalizeProcessName(match[1]),
			RawName: match[1],
			CPUTime: parseCPUTime(match[2]),
			MemPeak: parseMemPeak(match[3]),
			RawCPU:  match[2],
			RawMem:  match[3],
		})
	}

	return processes, scanner.Err()
}

func LoadThermalStats(since, until time.Time) (domain.ThermalStats, error) {
	args := []string{"--no-pager", "-o", "short-iso"}
	if !since.IsZero() {
		args = append(args, "--since", since.Format("2006-01-02"))
	}
	if !until.IsZero() {
		args = append(args, "--until", until.Format("2006-01-02 23:59:59"))
	}

	out, err := exec.Command("journalctl", args...).CombinedOutput()
	if err != nil && len(strings.TrimSpace(string(out))) == 0 {
		return domain.ThermalStats{}, err
	}

	var stats domain.ThermalStats
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

func LoadSpecs() domain.HardwareSpecs {
	specs := domain.HardwareSpecs{
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
	if out, err := exec.Command("sh", "-c", batCmd).CombinedOutput(); err == nil || len(strings.TrimSpace(string(out))) > 0 {
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

	if specs.Battery == "--" {
		specs.Battery = readBatterySysfs()
	}

	return specs
}

func readBatterySysfs() string {
	ueventPaths, _ := filepath.Glob("/sys/class/power_supply/BAT*/uevent")
	for _, path := range ueventPaths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var full, design, capPct, cycle string
		for _, line := range strings.Split(string(data), "\n") {
			switch {
			case strings.HasPrefix(line, "POWER_SUPPLY_ENERGY_FULL="):
				full = strings.TrimPrefix(line, "POWER_SUPPLY_ENERGY_FULL=")
			case strings.HasPrefix(line, "POWER_SUPPLY_ENERGY_FULL_DESIGN="):
				design = strings.TrimPrefix(line, "POWER_SUPPLY_ENERGY_FULL_DESIGN=")
			case strings.HasPrefix(line, "POWER_SUPPLY_CAPACITY="):
				capPct = strings.TrimPrefix(line, "POWER_SUPPLY_CAPACITY=")
			case strings.HasPrefix(line, "POWER_SUPPLY_CYCLE_COUNT="):
				cycle = strings.TrimPrefix(line, "POWER_SUPPLY_CYCLE_COUNT=")
			}
		}
		parts := []string{}
		if full != "" {
			if wh, err := strconv.ParseFloat(full, 64); err == nil {
				parts = append(parts, fmt.Sprintf("full %.1fWh", wh/1000000))
			}
		}
		if design != "" {
			if wh, err := strconv.ParseFloat(design, 64); err == nil {
				parts = append(parts, fmt.Sprintf("design %.1fWh", wh/1000000))
			}
		}
		if capPct != "" {
			parts = append(parts, capPct+"%")
		}
		if cycle != "" {
			parts = append(parts, "cycle "+cycle)
		}
		if len(parts) > 0 {
			return strings.Join(parts, ", ")
		}
	}
	return "--"
}

func parseHistoryFile(path string, since, until time.Time) ([]domain.BatteryPoint, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var points []domain.BatteryPoint
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
		points = append(points, domain.BatteryPoint{
			Time:       t,
			Percentage: pct,
			State:      parts[2],
		})
	}
	return points, scanner.Err()
}

func parseRateFile(path string, since, until time.Time) ([]domain.RatePoint, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var points []domain.RatePoint
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
		watts, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			continue
		}
		t := time.Unix(epoch, 0)
		if !since.IsZero() && t.Before(since) {
			continue
		}
		if !until.IsZero() && t.After(until) {
			continue
		}
		points = append(points, domain.RatePoint{
			Time:  t,
			Watts: watts,
			State: parts[2],
		})
	}
	return points, scanner.Err()
}

func filepathGlob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

func normalizeProcessName(raw string) string {
	name := decodeHexEscape(raw)
	name = strings.TrimSuffix(name, ".service")
	name = strings.TrimSuffix(name, ".scope")
	name = strings.TrimSuffix(name, ".slice")
	if strings.HasPrefix(name, "dbus-:") {
		if parts := strings.SplitN(name, "-", 3); len(parts) == 3 {
			name = parts[2]
		}
	}
	name = strings.TrimPrefix(name, "app-")
	name = strings.TrimPrefix(name, "flatpak-")
	name = strings.TrimPrefix(name, "gnome-")
	if idx := lastHyphenBeforeDigits(name); idx > 0 {
		name = name[:idx]
	}
	if idx := findUUIDStart(name); idx > 0 {
		name = name[:idx]
	}
	return strings.ToLower(strings.TrimSpace(name))
}

func decodeHexEscape(s string) string {
	var out strings.Builder
	for i := 0; i < len(s); {
		if i+3 < len(s) && s[i] == '\\' && s[i+1] == 'x' {
			if b, err := strconv.ParseUint(s[i+2:i+4], 16, 8); err == nil {
				out.WriteByte(byte(b))
				i += 4
				continue
			}
		}
		out.WriteByte(s[i])
		i++
	}
	return out.String()
}

func findUUIDStart(s string) int {
	for i := 0; i+9 < len(s); i++ {
		if s[i] != '-' {
			continue
		}
		if isHexGroup(s[i+1:], 8) && len(s) > i+9 && s[i+9] == '-' {
			return i
		}
	}
	return -1
}

func isHexGroup(s string, n int) bool {
	if len(s) < n {
		return false
	}
	for i := 0; i < n; i++ {
		c := s[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

func lastHyphenBeforeDigits(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] != '-' {
			continue
		}
		if i+1 >= len(s) {
			return -1
		}
		allDigits := true
		for _, c := range s[i+1:] {
			if c < '0' || c > '9' {
				allDigits = false
				break
			}
		}
		if allDigits {
			return i
		}
		return -1
	}
	return -1
}

func parseCPUTime(s string) float64 {
	var total float64
	fields := strings.Fields(strings.TrimSpace(s))
	for _, field := range fields {
		switch {
		case strings.HasSuffix(field, "h"):
			v, _ := strconv.ParseFloat(strings.TrimSuffix(field, "h"), 64)
			total += v * 3600
		case strings.HasSuffix(field, "min"):
			v, _ := strconv.ParseFloat(strings.TrimSuffix(field, "min"), 64)
			total += v * 60
		case strings.HasSuffix(field, "s"):
			v, _ := strconv.ParseFloat(strings.TrimSuffix(field, "s"), 64)
			total += v
		}
	}
	return total
}

func parseMemPeak(s string) float64 {
	s = strings.TrimSpace(strings.ToUpper(s))
	switch {
	case strings.HasSuffix(s, "G"):
		v, _ := strconv.ParseFloat(strings.TrimSuffix(s, "G"), 64)
		return v * 1024
	case strings.HasSuffix(s, "M"):
		v, _ := strconv.ParseFloat(strings.TrimSuffix(s, "M"), 64)
		return v
	case strings.HasSuffix(s, "K"):
		v, _ := strconv.ParseFloat(strings.TrimSuffix(s, "K"), 64)
		return v / 1024
	default:
		v, _ := strconv.ParseFloat(s, 64)
		return v / 1024 / 1024
	}
}
