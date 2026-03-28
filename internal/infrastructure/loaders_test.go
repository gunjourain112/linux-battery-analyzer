package infrastructure

import (
	"testing"
	"time"
)

func TestParseProcessUsageLine(t *testing.T) {
	line := "2026-03-03T14:49:05+09:00 fedora systemd[2017]: org.gnome.Shell@wayland.service: Consumed 1.944s CPU time, 429.6M memory peak."
	ts, ok := parseJournalTime(line)
	if !ok {
		t.Fatalf("parseJournalTime failed")
	}

	proc, ok := parseProcessUsageLine(line, ts)
	if !ok {
		t.Fatalf("parseProcessUsageLine failed")
	}
	if proc.RawName != "org.gnome.Shell@wayland.service" {
		t.Fatalf("unexpected RawName: %q", proc.RawName)
	}
	if proc.Name != "org.gnome.shell@wayland" {
		t.Fatalf("unexpected Name: %q", proc.Name)
	}
	if proc.CPUTime <= 1.9 || proc.CPUTime >= 2.0 {
		t.Fatalf("unexpected CPUTime: %v", proc.CPUTime)
	}
	if proc.MemPeak < 429 || proc.MemPeak > 430 {
		t.Fatalf("unexpected MemPeak: %v", proc.MemPeak)
	}
}

func TestParseThermalLine(t *testing.T) {
	line := "2026-03-03T23:48:27+09:00 fedora kernel: ACPI: thermal: Thermal Zone [THZ0] (39 C)"
	ts, ok := parseJournalTime(line)
	if !ok {
		t.Fatalf("parseJournalTime failed")
	}
	if got := ts.In(time.Local); got.IsZero() {
		t.Fatalf("parsed zero time")
	}

	temp, ok := parseThermalLine(line)
	if !ok {
		t.Fatalf("parseThermalLine failed")
	}
	if temp != 39 {
		t.Fatalf("unexpected temp: %d", temp)
	}
}
