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

	case HighlightTextStyleID:
		s = s.Foreground(d.Color(HighlightFontColorID))

	case DimTextStyleID:
		s = s.Foreground(d.Color(DimFontColorID))

	case DimSecondaryTextStyleID:
		s = s.Italic(true).Foreground(d.Color(DimFontColorID))

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

	case StatusErrorTextStyleID:
		s = s.AlignHorizontal(lipgloss.Center).
			Bold(true).Foreground(d.Color(StatusErrorFontColorID))

	case StatusSuccessTextStyleID:
		s = s.AlignHorizontal(lipgloss.Center).
			Bold(true).Foreground(d.Color(StatusSuccessFontColorID))

	case StatusWarningTextStyleID:
		s = s.AlignHorizontal(lipgloss.Center).
			Foreground(d.Color(StatusWarningFontColorID))

	case StatusInformationTextStyleID:
		s = s.AlignHorizontal(lipgloss.Center).
			Foreground(d.Color(StatusInformationFontColorID))

	case StatusNeutralTextStyleID:
		s = s.AlignHorizontal(lipgloss.Center).
			Foreground(d.Color(StatusNeutralFontColorID))

	}

	return s
}

func (d DefaultDarkTheme) Color(color ColorID) lipgloss.Color {
	var colorHexCode string

	switch color {
	case NeutralFontColorID:
		colorHexCode = "#F5F5F5"

	case BlurredBorderColorID, BlurredFontColorID, DimFontColorID:
		colorHexCode = "#186318"

	case HighlightFontColorID:
		colorHexCode = "#C7FF37"

	case StatusErrorFontColorID:
		colorHexCode = "#DB0000"

	case StatusSuccessFontColorID:
		colorHexCode = "#27DB18"

	case StatusWarningFontColorID:
		colorHexCode = "#FF7B00"

	case StatusInformationFontColorID:
		colorHexCode = "#039FFC"

	case StatusNeutralFontColorID:
		colorHexCode = "#D0D0D0"

	default:
		colorHexCode = "#18B718"

	}

	return lipgloss.Color(colorHexCode)
}

func (d DefaultDarkTheme) Size(size SizeID) int {
	return 0
}
