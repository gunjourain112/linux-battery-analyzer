package domain

import "time"

type BatteryPoint struct {
	Time       time.Time
	Percentage float64
	State      string
}

type RatePoint struct {
	Time  time.Time
	Watts float64
	State string
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

type LoadLevel int

const (
	LoadLevelUnknown LoadLevel = iota
	LoadLevelLight
	LoadLevelMedium
	LoadLevelHeavy
)

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

type LoadBucket struct {
	Label    string
	Count    int
	Ratio    float64
	AvgWatts float64
	EstHours float64
}

type DischargeProfile struct {
	Buckets    []LoadBucket
	TotalCount int
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

type ProcessImpact struct {
	Process    ProcessUsage
	DrainWatts float64
	Level      LoadLevel
}

type ChargingSession struct {
	Start       time.Time
	End         time.Time
	StartPct    float64
	EndPct      float64
	AvgChargeW  float64
	PeakChargeW float64
}

type SystemEvent struct {
	Time time.Time
	Type string
	Desc string
}

type DailyRecord struct {
	Date      time.Time
	ActiveMin int
	Discharge float64
	Charge    float64
	AvgWatts  float64
}
