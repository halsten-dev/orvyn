package theme

import (
	"github.com/charmbracelet/lipgloss"
)

type DefaultDarkTheme struct {
	theme Theme
}

func NewDefaultDarkTheme() *DefaultDarkTheme {
	d := &DefaultDarkTheme{}
	d.theme = d

	return d
}

func (d DefaultDarkTheme) Style(style StyleID) lipgloss.Style {
	var s lipgloss.Style

	s = lipgloss.NewStyle()

	switch style {
	case TitleStyleID:
		s = s.Bold(true).Foreground(d.theme.Color(TitleFontColorID))

	case NeutralTextStyleID:
		s = s.Foreground(d.theme.Color(NeutralFontColorID))

	case NeutralDimTextStyleID:
		s = s.Foreground(d.theme.Color(NeutralDimFontColorID))

	case NormalTextStyleID:
		s = s.Foreground(d.theme.Color(NormalFontColorID))

	case HighlightTextStyleID:
		s = s.Foreground(d.theme.Color(HighlightFontColorID))

	case DimTextStyleID:
		s = s.Foreground(d.theme.Color(DimFontColorID))

	case DimSecondaryTextStyleID:
		s = s.Italic(true).Foreground(d.theme.Color(DimFontColorID))

	case FocusedWidgetStyleID:
		s = s.Border(lipgloss.RoundedBorder()).
			BorderForeground(d.theme.Color(FocusedBorderColorID))

	case BlurredWidgetStyleID:
		s = s.Border(lipgloss.RoundedBorder()).
			BorderForeground(d.theme.Color(BlurredBorderColorID))

	case PaginatorActiveStyleID:
		s = s.Bold(true).Foreground(d.theme.Color(NormalFontColorID))

	case PaginatorInactiveStyleID:
		s = s.Foreground(d.theme.Color(DimFontColorID))

	case StatusErrorTextStyleID:
		s = s.AlignHorizontal(lipgloss.Center).
			Bold(true).Foreground(d.theme.Color(StatusErrorFontColorID))

	case StatusSuccessTextStyleID:
		s = s.AlignHorizontal(lipgloss.Center).
			Bold(true).Foreground(d.theme.Color(StatusSuccessFontColorID))

	case StatusWarningTextStyleID:
		s = s.AlignHorizontal(lipgloss.Center).
			Foreground(d.theme.Color(StatusWarningFontColorID))

	case StatusInformationTextStyleID:
		s = s.AlignHorizontal(lipgloss.Center).
			Foreground(d.theme.Color(StatusInformationFontColorID))

	case StatusNeutralTextStyleID:
		s = s.AlignHorizontal(lipgloss.Center).
			Foreground(d.theme.Color(StatusNeutralFontColorID))

	}

	return s
}

func (d DefaultDarkTheme) Color(color ColorID) lipgloss.Color {
	var colorHexCode string

	switch color {
	case NeutralFontColorID:
		colorHexCode = "#F5F5F5"

	case NeutralDimFontColorID:
		colorHexCode = "#898989"

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
