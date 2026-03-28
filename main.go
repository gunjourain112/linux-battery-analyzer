package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/infrastructure"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/service"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/renderer"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/tui"
)

func main() {
	conf, err := loadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	points, err := infrastructure.LoadHistory(conf.Since, conf.Until)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load history:", err)
		os.Exit(1)
	}

	events, err := infrastructure.LoadPowerEvents(conf.Since, conf.Until)
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning: could not load journal:", err)
	}

	processes, err := infrastructure.LoadProcessUsage(conf.Since, conf.Until)
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning: could not load resource logs:", err)
	}

	rates, err := infrastructure.LoadRateHistory(conf.Since, conf.Until)
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning: could not load rate logs:", err)
	}

	thermal, err := infrastructure.LoadThermalStats(conf.Since, conf.Until)
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning: could not load thermal logs:", err)
	}

	specs := infrastructure.LoadSpecs()

	sessions := service.BuildSessions(events, points, conf.Since, conf.Until)
	profile := service.BuildDischargeProfile(rates)
	impacts := service.BuildProcessImpacts(processes, rates)
	charging := service.BuildChargingSessions(points, rates)
	daily := service.BuildDailySummary(sessions, charging)
	systemEvents := service.BuildSystemEvents(events, points)
	detailed := service.BuildDetailedTimeline(points, rates, systemEvents, processes, thermal)

	report := renderer.ReportData{
		Config:         conf,
		Sessions:       sessions,
		Detailed:       detailed,
		BatteryHistory: points,
		Charging:       charging,
		Daily:          daily,
		SystemEvents:   systemEvents,
		Discharge:      profile,
		ProcessImpacts: impacts,
		Processes:      processes,
		Specs:          specs,
		Thermal:        thermal,
	}

	fmt.Print(renderer.Render(report))
}

func loadConfig() (domain.Config, error) {
	if len(os.Args) >= 3 {
		since, err := time.ParseInLocation("2006-01-02", os.Args[1], time.Local)
		if err != nil {
			return domain.Config{}, fmt.Errorf("invalid since date: %s", os.Args[1])
		}
		until, err := time.ParseInLocation("2006-01-02", os.Args[2], time.Local)
		if err != nil {
			return domain.Config{}, fmt.Errorf("invalid until date: %s", os.Args[2])
		}
		return domain.Config{
			Language: "ko",
			Since:    since,
			Until:    until.Add(24*time.Hour - time.Second),
		}, nil
	}

	return tui.Run()
}
