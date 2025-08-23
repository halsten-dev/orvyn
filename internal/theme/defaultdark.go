package theme

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/halsten-dev/orvyn"
)

type DefaultDarkTheme struct{}

func (d *DefaultDarkTheme) Style(style orvyn.StyleName) lipgloss.Style {
	var s lipgloss.Style

	s = lipgloss.NewStyle()

	switch style {
	case orvyn.TitleStyleName:
		s = s.Bold(true).Foreground(d.Color(orvyn.TitleFontColorName))

	case orvyn.NormalTextName:
		s = s.Foreground(d.Color(orvyn.NormalFontColorName))

	case orvyn.FocusedWidgetStyleName:
		s = s.Border(lipgloss.RoundedBorder()).
			BorderForeground(d.Color(orvyn.FocusedBorderColorName))

	case orvyn.BlurredWidgetStyleName:
		s = s.Border(lipgloss.RoundedBorder()).
			BorderForeground(d.Color(orvyn.BlurredBorderColorName))

	}

	return s
}

func (d *DefaultDarkTheme) Color(color orvyn.ColorName) lipgloss.Color {
	var colorHexCode string

	switch color {
	case orvyn.BlurredBorderColorName, orvyn.BlurredFontColorName:
		colorHexCode = "#186318"

	default:
		colorHexCode = "#18B718"

	}

	return lipgloss.Color(colorHexCode)
}

func (d *DefaultDarkTheme) Size(size orvyn.SizeName) int {
	return 0
}
