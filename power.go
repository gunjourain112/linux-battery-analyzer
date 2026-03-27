package main

import (
	"bufio"
	"os/exec"
	"strings"
	"time"
)

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
