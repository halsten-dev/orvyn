package theme

import (
	"github.com/charmbracelet/lipgloss"
)

type DefaultDarkTheme struct{}

func (d DefaultDarkTheme) Style(style StyleName) lipgloss.Style {
	var s lipgloss.Style

	s = lipgloss.NewStyle()

	switch style {
	case TitleStyleName:
		s = s.Bold(true).Foreground(d.Color(TitleFontColorName))

	case NormalTextName:
		s = s.Foreground(d.Color(NormalFontColorName))

	case FocusedWidgetStyleName:
		s = s.Border(lipgloss.RoundedBorder()).
			BorderForeground(d.Color(FocusedBorderColorName))

	case BlurredWidgetStyleName:
		s = s.Border(lipgloss.RoundedBorder()).
			BorderForeground(d.Color(BlurredBorderColorName))

	}

	return s
}

func (d DefaultDarkTheme) Color(color ColorName) lipgloss.Color {
	var colorHexCode string

	switch color {
	case BlurredBorderColorName, BlurredFontColorName:
		colorHexCode = "#186318"

	default:
		colorHexCode = "#18B718"

	}

	return lipgloss.Color(colorHexCode)
}

func (d DefaultDarkTheme) Size(size SizeName) int {
	return 0
}
