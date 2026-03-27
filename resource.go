package main

import (
	"bufio"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var resourceLine = regexp.MustCompile(`: (?:app-)?(?:gnome-|flatpak-)?(.*)\.scope: Consumed (.*) CPU time, (.*) memory peak\.`)

func loadProcessUsage(since, until time.Time) ([]ProcessUsage, error) {
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

	var processes []ProcessUsage
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

		processes = append(processes, ProcessUsage{
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

func normalizeProcessName(raw string) string {
	name := decodeHexEscape(raw)
	if idx := strings.LastIndex(name, ".scope"); idx >= 0 {
		name = strings.TrimSuffix(name, ".scope")
	}
	if idx := lastHyphenBeforeDigits(name); idx > 0 {
		name = name[:idx]
	}
	if idx := findUUIDStart(name); idx > 0 {
		name = name[:idx]
	}
	if parts := strings.Split(name, "."); len(parts) > 1 {
		name = parts[len(parts)-1]
	}
	name = strings.TrimPrefix(name, "app-")
	name = strings.TrimPrefix(name, "flatpak-")
	name = strings.TrimPrefix(name, "gnome-")
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
