package renderer

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
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

	title := "Notebook Battery Analyzer"
	if d.Config.Language == "en" {
		title = "Notebook Battery Analyzer"
	}

	b.WriteString(t.Title().Render(title))
	b.WriteString("\n")
	if !d.Config.Since.IsZero() || !d.Config.Until.IsZero() {
		b.WriteString(t.Subtle().Render(renderRange(d.Config.Since, d.Config.Until)))
		b.WriteString("\n")
	}
	if !d.Specs.IsEmpty() {
		b.WriteString(t.Subtle().Render(renderSpecsLine(d.Specs)))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	b.WriteString(renderSection("Summary", renderSummary(d)))
	b.WriteString(renderSection("Sessions", renderSessions(d.Sessions)))
	b.WriteString(renderSection("Daily", renderDaily(d.Daily)))
	b.WriteString(renderSection("Charging", renderCharging(d.Charging)))
	b.WriteString(renderSection("Discharge Profile", renderProfile(d.Discharge)))
	b.WriteString(renderSection("Process Impacts", renderImpacts(d.ProcessImpacts)))
	b.WriteString(renderSection("System Events", renderEvents(d.SystemEvents)))
	b.WriteString(renderSection("Specs", renderSpecsTable(d.Specs)))
	b.WriteString(renderSection("Thermals", renderThermals(d.Thermal)))

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

func renderSummary(d ReportData) string {
	var totalHours float64
	var totalDrain float64
	var worst domain.Session
	var worstRate float64
	var hasWorst bool

	for _, s := range d.Sessions {
		rate := dischargeRate(s)
		if rate <= 0 {
			continue
		}
		totalHours += s.End.Sub(s.Start).Hours()
		totalDrain += s.StartPct - s.EndPct
		if !hasWorst || rate > worstRate {
			hasWorst = true
			worst = s
			worstRate = rate
		}
	}

	chargingCount := len(d.Charging)
	eventsCount := len(d.SystemEvents)

	lines := []string{
		fmt.Sprintf("sessions: %d", len(d.Sessions)),
		fmt.Sprintf("charging: %d", chargingCount),
		fmt.Sprintf("events: %d", eventsCount),
	}
	if totalHours > 0 {
		lines = append(lines, fmt.Sprintf("avg discharge: %.2f%%/h", totalDrain/totalHours))
	} else {
		lines = append(lines, "avg discharge: --")
	}
	if hasWorst {
		lines = append(lines, fmt.Sprintf("worst session: %s ~ %s (%.2f%%/h)",
			worst.Start.Format("01/02 15:04"),
			worst.End.Format("15:04"),
			worstRate,
		))
	}

	return joinLines(lines...)
}

func renderSessions(sessions []domain.Session) string {
	if len(sessions) == 0 {
		return theme.Default.Subtle().Render("no sessions")
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers("Start", "End", "Duration", "Drain", "Rate").
		StyleFunc(func(r, c int) lipgloss.Style {
			if r == -1 {
				return theme.Default.Header()
			}
			if c >= 2 {
				return theme.Default.Value().Align(lipgloss.Right)
			}
			return theme.Default.Value()
		})

	for _, s := range sessions {
		rate := dischargeRate(s)
		drain := s.StartPct - s.EndPct
		tbl.Row(
			s.Start.Format("01/02 15:04"),
			s.End.Format("01/02 15:04"),
			formatDuration(s.End.Sub(s.Start)),
			fmt.Sprintf("%.0f%%", drain),
			fmt.Sprintf("%.2f%%/h", rate),
		)
	}
	return tbl.Render()
}

func renderDaily(records []domain.DailyRecord) string {
	if len(records) == 0 {
		return theme.Default.Subtle().Render("no daily records")
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers("Date", "Active", "Discharge", "Charge", "Avg W").
		StyleFunc(func(r, c int) lipgloss.Style {
			if r == -1 {
				return theme.Default.Header()
			}
			if c == 0 {
				return theme.Default.Value()
			}
			return theme.Default.Value().Align(lipgloss.Right)
		})

	for _, r := range records {
		tbl.Row(
			r.Date.Format("01/02"),
			fmt.Sprintf("%dm", r.ActiveMin),
			fmt.Sprintf("%.0f%%", r.Discharge),
			fmt.Sprintf("%.0f%%", r.Charge),
			fmt.Sprintf("%.1f", r.AvgWatts),
		)
	}
	return tbl.Render()
}

func renderCharging(sessions []domain.ChargingSession) string {
	if len(sessions) == 0 {
		return theme.Default.Subtle().Render("no charging sessions")
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers("Start", "End", "Duration", "Gain", "Avg W", "Peak W").
		StyleFunc(func(r, c int) lipgloss.Style {
			if r == -1 {
				return theme.Default.Header()
			}
			if c >= 2 {
				return theme.Default.Value().Align(lipgloss.Right)
			}
			return theme.Default.Value()
		})

	for _, s := range sessions {
		gain := s.EndPct - s.StartPct
		tbl.Row(
			s.Start.Format("01/02 15:04"),
			s.End.Format("01/02 15:04"),
			formatDuration(s.End.Sub(s.Start)),
			fmt.Sprintf("%.0f%%", gain),
			fmt.Sprintf("%.1f", s.AvgChargeW),
			fmt.Sprintf("%.1f", s.PeakChargeW),
		)
	}
	return tbl.Render()
}

func renderProfile(profile domain.DischargeProfile) string {
	if profile.TotalCount == 0 {
		return theme.Default.Subtle().Render("no discharge profile")
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers("Bucket", "Count", "Ratio", "Avg W").
		StyleFunc(func(r, c int) lipgloss.Style {
			if r == -1 {
				return theme.Default.Header()
			}
			if c == 0 {
				return theme.Default.Label()
			}
			return theme.Default.Value().Align(lipgloss.Right)
		})

	for _, b := range profile.Buckets {
		if b.Count == 0 {
			continue
		}
		tbl.Row(
			b.Label,
			fmt.Sprintf("%d", b.Count),
			fmt.Sprintf("%.0f%%", b.Ratio),
			fmt.Sprintf("%.1f", b.AvgWatts),
		)
	}
	return tbl.Render()
}

func renderImpacts(impacts []domain.ProcessImpact) string {
	if len(impacts) == 0 {
		return theme.Default.Subtle().Render("no process impact data")
	}

	limit := len(impacts)
	if limit > 5 {
		limit = 5
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers("Process", "Drain W", "Level", "CPU s", "Mem M").
		StyleFunc(func(r, c int) lipgloss.Style {
			if r == -1 {
				return theme.Default.Header()
			}
			if c == 0 {
				return theme.Default.Value()
			}
			return theme.Default.Value().Align(lipgloss.Right)
		})

	for i := 0; i < limit; i++ {
		p := impacts[i]
		tbl.Row(
			p.Process.Name,
			fmt.Sprintf("%.1f", p.DrainWatts),
			levelString(p.Level),
			fmt.Sprintf("%.0f", p.Process.CPUTime),
			fmt.Sprintf("%.1f", p.Process.MemPeak),
		)
	}
	return tbl.Render()
}

func renderEvents(events []domain.SystemEvent) string {
	if len(events) == 0 {
		return theme.Default.Subtle().Render("no system events")
	}

	limit := len(events)
	if limit > 8 {
		limit = 8
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers("Time", "Type", "Description").
		StyleFunc(func(r, c int) lipgloss.Style {
			if r == -1 {
				return theme.Default.Header()
			}
			if c == 0 {
				return theme.Default.Value()
			}
			return theme.Default.Value()
		})

	for i := 0; i < limit; i++ {
		ev := events[i]
		tbl.Row(
			ev.Time.Format("01/02 15:04"),
			ev.Type,
			ev.Desc,
		)
	}
	return tbl.Render()
}

func renderSpecsTable(specs domain.HardwareSpecs) string {
	if specs.IsEmpty() {
		return theme.Default.Subtle().Render("no hardware specs")
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		StyleFunc(func(r, c int) lipgloss.Style {
			if c == 0 {
				return theme.Default.Label()
			}
			return theme.Default.Value()
		})

	if specs.OS != "" {
		tbl.Row("OS", specs.OS)
	}
	if specs.Device != "" {
		tbl.Row("Device", specs.Device)
	}
	if specs.CPU != "" {
		tbl.Row("CPU", specs.CPU)
	}
	if specs.RAM != "" {
		tbl.Row("RAM", specs.RAM)
	}
	if specs.Battery != "" {
		tbl.Row("Battery", specs.Battery)
	}

	return tbl.Render()
}

func renderThermals(stats domain.ThermalStats) string {
	if stats.Count == 0 {
		return theme.Default.Subtle().Render("no thermal samples")
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(theme.Default.Subtle()).
		Headers("Samples", "Min", "Max", "Avg").
		StyleFunc(func(r, c int) lipgloss.Style {
			if r == -1 {
				return theme.Default.Header()
			}
			return theme.Default.Value().Align(lipgloss.Right)
		})

	tbl.Row(
		fmt.Sprintf("%d", stats.Count),
		fmt.Sprintf("%d C", stats.Min),
		fmt.Sprintf("%d C", stats.Max),
		fmt.Sprintf("%d C", stats.Avg),
	)

	return tbl.Render()
}

func renderRange(since, until time.Time) string {
	if since.IsZero() && until.IsZero() {
		return ""
	}
	if since.IsZero() {
		return fmt.Sprintf("until %s", until.Format("2006-01-02"))
	}
	if until.IsZero() {
		return fmt.Sprintf("since %s", since.Format("2006-01-02"))
	}
	return fmt.Sprintf("%s ~ %s", since.Format("2006-01-02"), until.Format("2006-01-02"))
}

func renderSpecsLine(specs domain.HardwareSpecs) string {
	parts := make([]string, 0, 4)
	if specs.Device != "" {
		parts = append(parts, specs.Device)
	}
	if specs.CPU != "" {
		parts = append(parts, specs.CPU)
	}
	if specs.RAM != "" {
		parts = append(parts, specs.RAM)
	}
	if specs.Battery != "" {
		parts = append(parts, specs.Battery)
	}
	return strings.Join(parts, "  |  ")
}

func joinLines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = -d
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh%02dm", h, m)
}

func dischargeRate(s domain.Session) float64 {
	dur := s.End.Sub(s.Start).Hours()
	if dur <= 0 {
		return 0
	}
	drain := s.StartPct - s.EndPct
	if drain <= 0 {
		return 0
	}
	return drain / dur
}

func levelString(level domain.LoadLevel) string {
	switch level {
	case domain.LoadLevelLight:
		return "light"
	case domain.LoadLevelMedium:
		return "medium"
	case domain.LoadLevelHeavy:
		return "heavy"
	default:
		return "unknown"
	}
}
