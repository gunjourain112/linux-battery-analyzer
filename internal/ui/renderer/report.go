package renderer

import (
	"strings"

	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/renderer/components"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

type ReportData struct {
	Config         domain.Config
	Sessions       []domain.Session
	Charging       []domain.ChargingSession
	Daily          []domain.DailyRecord
	SystemEvents   []domain.SystemEvent
	Discharge      domain.DischargeProfile
	ProcessImpacts []domain.ProcessImpact
	Processes      []domain.ProcessUsage
	Specs          domain.HardwareSpecs
	Thermal        domain.ThermalStats
}

func Render(d ReportData) string {
	t := theme.Default
	var b strings.Builder

	b.WriteString(t.Title().Render("Notebook Battery Analyzer"))
	b.WriteString("\n")
	if rangeLine := components.RenderRange(d.Config.Since, d.Config.Until); rangeLine != "" {
		b.WriteString(t.Subtle().Render(rangeLine))
		b.WriteString("\n")
	}
	if specsLine := components.RenderSpecsLine(d.Specs); specsLine != "" {
		b.WriteString(t.Subtle().Render(specsLine))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	b.WriteString(renderSection("Summary", components.Summary(d.Sessions, d.Charging, d.SystemEvents)))
	b.WriteString(renderSection("Sessions", components.Sessions(d.Sessions)))
	b.WriteString(renderSection("Daily", components.Daily(d.Daily)))
	b.WriteString(renderSection("Charging", components.Charging(d.Charging)))
	b.WriteString(renderSection("Discharge Profile", components.DischargeProfile(d.Discharge)))
	b.WriteString(renderSection("Process Impacts", components.ProcessImpacts(d.ProcessImpacts)))
	b.WriteString(renderSection("System Events", components.SystemEvents(d.SystemEvents)))
	b.WriteString(renderSection("Specs", components.Specs(d.Specs)))
	b.WriteString(renderSection("Thermals", components.Thermals(d.Thermal)))

	return b.String()
}

func renderSection(title, body string) string {
	if strings.TrimSpace(body) == "" {
		return ""
	}

	t := theme.Default
	var b strings.Builder
	b.WriteString(t.Subtle().Render(strings.Repeat("─", 72)))
	b.WriteString("\n")
	b.WriteString(t.SectionTitle().Render(" " + title + " "))
	b.WriteString("\n\n")
	b.WriteString(body)
	if !strings.HasSuffix(body, "\n") {
		b.WriteString("\n")
	}
	b.WriteString("\n")
	return b.String()
}
