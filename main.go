package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/infrastructure"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/service"
)

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

	points, err := infrastructure.LoadHistory(since, until)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load history:", err)
		os.Exit(1)
	}

	events, err := infrastructure.LoadPowerEvents(since, until)
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning: could not load journal:", err)
	}

	processes, err := infrastructure.LoadProcessUsage(since, until)
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning: could not load resource logs:", err)
	}

	thermal, err := infrastructure.LoadThermalStats(since, until)
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning: could not load thermal logs:", err)
	}

	specs := infrastructure.LoadSpecs()

	sessions := service.BuildSessions(events, points, since, until)

	fmt.Println("=== sessions ===")
	for i, s := range sessions {
		dur := s.End.Sub(s.Start)
		drain := s.StartPct - s.EndPct
		rate := service.DischargeRate(s)
		fmt.Printf("[%d] %s ~ %s  (%dh%02dm)  %.0f%% → %.0f%%  (%.1f%% drain)\n",
			i+1,
			s.Start.Format("01/02 15:04"),
			s.End.Format("15:04"),
			int(dur.Hours()), int(dur.Minutes())%60,
			s.StartPct, s.EndPct, drain,
		)
		if rate > 0 {
			fmt.Printf("    rate: %.2f%%/h\n", rate)
		}
	}

	printSummary(sessions)
	printSpecs(specs)
	printThermals(thermal)
	printProcesses(processes)
}

func printSummary(sessions []domain.Session) {
	var totalHours float64
	var totalDrain float64
	var worst domain.Session
	var worstRate float64
	var hasWorst bool

	for _, s := range sessions {
		rate := service.DischargeRate(s)
		if rate <= 0 {
			continue
		}

		hours := s.End.Sub(s.Start).Hours()
		totalHours += hours
		totalDrain += s.StartPct - s.EndPct

		if !hasWorst || rate > worstRate {
			hasWorst = true
			worst = s
			worstRate = rate
		}
	}

	fmt.Println()
	fmt.Println("=== summary ===")
	fmt.Printf("sessions: %d\n", len(sessions))

	if totalHours == 0 {
		fmt.Println("avg discharge: --")
		fmt.Println("worst session: --")
		return
	}

	avgRate := totalDrain / totalHours
	fmt.Printf("avg discharge: %.2f%%/h\n", avgRate)

	if hasWorst {
		dur := worst.End.Sub(worst.Start)
		fmt.Printf(
			"worst session: %s ~ %s  (%dh%02dm)  %.2f%%/h\n",
			worst.Start.Format("01/02 15:04"),
			worst.End.Format("15:04"),
			int(dur.Hours()), int(dur.Minutes())%60,
			worstRate,
		)
	} else {
		fmt.Println("worst session: --")
	}
}

func printProcesses(processes []domain.ProcessUsage) {
	if len(processes) == 0 {
		return
	}

	sort.Slice(processes, func(i, j int) bool {
		if processes[i].CPUTime == processes[j].CPUTime {
			return processes[i].MemPeak > processes[j].MemPeak
		}
		return processes[i].CPUTime > processes[j].CPUTime
	})

	limit := 5
	if len(processes) < limit {
		limit = len(processes)
	}

	fmt.Println()
	fmt.Println("=== processes ===")
	for i := 0; i < limit; i++ {
		p := processes[i]
		fmt.Printf("[%d] %s  cpu %.0fs  mem %.1fM\n", i+1, p.Name, p.CPUTime, p.MemPeak)
	}
}

func printSpecs(specs domain.HardwareSpecs) {
	if specs.IsEmpty() {
		return
	}

	fmt.Println()
	fmt.Println("=== specs ===")
	fmt.Printf("os: %s\n", specs.OS)
	fmt.Printf("device: %s\n", specs.Device)
	fmt.Printf("cpu: %s\n", specs.CPU)
	fmt.Printf("ram: %s\n", specs.RAM)
	fmt.Printf("battery: %s\n", specs.Battery)
}

func printThermals(stats domain.ThermalStats) {
	if stats.Count == 0 {
		return
	}

	fmt.Println()
	fmt.Println("=== thermal ===")
	fmt.Printf("samples: %d\n", stats.Count)
	fmt.Printf("min/max/avg: %d / %d / %d C\n", stats.Min, stats.Max, stats.Avg)
}
