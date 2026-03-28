package renderer

import (
	"strings"

	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/i18n"
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
	tr := i18n.New(d.Config.Language)
	var b strings.Builder

	b.WriteString(t.Title().Render(tr.Get(i18n.AppTitle)))
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

	b.WriteString(renderSection(tr.Get(i18n.ReportSummary), components.Summary(tr, d.Sessions, d.Charging, d.SystemEvents)))
	b.WriteString(renderSection(tr.Get(i18n.ReportSessions), components.Sessions(tr, d.Sessions)))
	b.WriteString(renderSection(tr.Get(i18n.ReportDaily), components.Daily(tr, d.Daily)))
	b.WriteString(renderSection(tr.Get(i18n.ReportCharging), components.Charging(tr, d.Charging)))
	b.WriteString(renderSection(tr.Get(i18n.ReportDischargeProfile), components.DischargeProfile(tr, d.Discharge)))
	b.WriteString(renderSection(tr.Get(i18n.ReportProcessImpacts), components.ProcessImpacts(tr, d.ProcessImpacts)))
	b.WriteString(renderSection(tr.Get(i18n.ReportSystemEvents), components.SystemEvents(tr, d.SystemEvents)))
	b.WriteString(renderSection(tr.Get(i18n.ReportSpecs), components.Specs(tr, d.Specs)))
	b.WriteString(renderSection(tr.Get(i18n.ReportThermals), components.Thermals(tr, d.Thermal)))

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
