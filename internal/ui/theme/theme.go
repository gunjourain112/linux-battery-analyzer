package theme

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Accent    lipgloss.Color
	Gray      lipgloss.Color
	Success   lipgloss.Color
	Warning   lipgloss.Color
	Error     lipgloss.Color
}

var Default = Theme{
	Primary:   lipgloss.Color("81"),
	Secondary: lipgloss.Color("212"),
	Accent:    lipgloss.Color("111"),
	Gray:      lipgloss.Color("246"),
	Success:   lipgloss.Color("70"),
	Warning:   lipgloss.Color("214"),
	Error:     lipgloss.Color("203"),
}

func (t Theme) Title() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true)
}

func (t Theme) SectionTitle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Background(t.Secondary).
		Padding(0, 1).
		Bold(true)
}

func (t Theme) Subtle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Gray)
}

func (t Theme) Label() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true)
}

func (t Theme) Value() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
}

func (t Theme) Header() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Bold(true)
}

func (t Theme) Good() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Success).Bold(true)
}

func (t Theme) WarningText() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Warning).Bold(true)
}

func (t Theme) ErrorText() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Error).Bold(true)
}
