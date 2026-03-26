package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	bat, err := getBatteryDevice()
	if err != nil {
		fmt.Fprintln(os.Stderr, "upower device not found:", err)
		os.Exit(1)
	}

	out, err := exec.Command("upower", "-i", bat).Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "upower failed:", err)
		os.Exit(1)
	}

	fmt.Print(string(out))
}

func getBatteryDevice() (string, error) {
	out, err := exec.Command("upower", "-e").Output()
	if err != nil {
		return "", err
	}
	for _, line := range splitLines(string(out)) {
		if len(line) > 0 && containsBAT(line) {
			return line, nil
		}
	}
	return "", fmt.Errorf("no battery found")
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	return lines
}

func containsBAT(s string) bool {
	for i := 0; i+3 <= len(s); i++ {
		if s[i] == 'B' && s[i+1] == 'A' && s[i+2] == 'T' {
			return true
		}
	}
	return false
}
