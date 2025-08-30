package theme

import (
	"github.com/charmbracelet/lipgloss"
)

type DefaultDarkTheme struct{}

func (d DefaultDarkTheme) Style(style StyleID) lipgloss.Style {
	var s lipgloss.Style

	s = lipgloss.NewStyle()

	switch style {
	case TitleStyleID:
		s = s.Bold(true).Foreground(d.Color(TitleFontColorID))

	case NormalTextStyleID:
		s = s.Foreground(d.Color(NormalFontColorID))

	case FocusedWidgetStyleID:
		s = s.Border(lipgloss.RoundedBorder()).
			BorderForeground(d.Color(FocusedBorderColorID))

	case BlurredWidgetStyleID:
		s = s.Border(lipgloss.RoundedBorder()).
			BorderForeground(d.Color(BlurredBorderColorID))

	case PaginatorActiveStyleID:
		s = s.Bold(true).Foreground(d.Color(NormalFontColorID))

	case PaginatorInactiveStyleID:
		s = s.Foreground(d.Color(DimFontColorID))
	}

	return s
}

func (d DefaultDarkTheme) Color(color ColorID) lipgloss.Color {
	var colorHexCode string

	switch color {
	case BlurredBorderColorID, BlurredFontColorID, DimFontColorID:
		colorHexCode = "#186318"

	default:
		colorHexCode = "#18B718"

	}

	return lipgloss.Color(colorHexCode)
}

func (d DefaultDarkTheme) Size(size SizeID) int {
	return 0
}
