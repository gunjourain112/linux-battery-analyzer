package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/i18n"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/theme"
)

func Specs(tr i18n.Translator, specs domain.HardwareSpecs) string {
	if specs.IsEmpty() {
		return theme.Default.Subtle().Render(tr.Get(i18n.NoHardwareSpecs))
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
		tbl.Row(tr.Get(i18n.OSHeader), specs.OS)
	}
	if specs.Device != "" {
		tbl.Row(tr.Get(i18n.DeviceHeader), specs.Device)
	}
	if specs.CPU != "" {
		tbl.Row(tr.Get(i18n.CPUHeader), specs.CPU)
	}
	if specs.RAM != "" {
		tbl.Row(tr.Get(i18n.RAMHeader), specs.RAM)
	}
	if specs.Battery != "" {
		tbl.Row(tr.Get(i18n.BatteryHeader), specs.Battery)
	}

	return tbl.Render()
}
