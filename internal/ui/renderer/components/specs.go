package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func Specs(specs domain.HardwareSpecs) string {
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
