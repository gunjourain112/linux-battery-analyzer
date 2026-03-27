package domain

import "time"

type BatteryPoint struct {
	Time       time.Time
	Percentage float64
	State      string
}

type PowerEvent struct {
	Time time.Time
	Kind string
}

type Session struct {
	Start    time.Time
	End      time.Time
	StartPct float64
	EndPct   float64
}

type ProcessUsage struct {
	Time    time.Time
	Name    string
	RawName string
	CPUTime float64
	MemPeak float64
	RawCPU  string
	RawMem  string
}

type ThermalStats struct {
	Count int
	Min   int
	Max   int
	Avg   int
}

type HardwareSpecs struct {
	OS      string
	Device  string
	CPU     string
	RAM     string
	Battery string
}

func (s HardwareSpecs) IsEmpty() bool {
	return s.OS == "" && s.Device == "" && s.CPU == "" && s.RAM == "" && s.Battery == ""
}
