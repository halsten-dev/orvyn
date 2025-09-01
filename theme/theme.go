package theme

import (
	"github.com/charmbracelet/lipgloss"
)

type StyleID uint

const (
	TitleStyleID StyleID = iota
	NormalTextStyleID
	HighlightTextStyleID
	DimTextStyleID
	DimSecondaryTextStyleID
	LabelTextStyleID
	FocusedWidgetStyleID
	BlurredWidgetStyleID
	PaginatorActiveStyleID
	PaginatorInactiveStyleID
	StatusErrorTextStyleID
	StatusSuccessTextStyleID
	StatusWarningTextStyleID
	StatusInformationTextStyleID
	StatusNeutralTextStyleID
)

type ColorID uint

const (
	TitleFontColorID ColorID = iota
	NormalFontColorID
	HighlightFontColorID
	DimFontColorID
	FocusedBorderColorID
	FocusedFontColorID
	BlurredBorderColorID
	BlurredFontColorID
	StatusErrorFontColorID
	StatusSuccessFontColorID
	StatusWarningFontColorID
	StatusInformationFontColorID
	StatusNeutralFontColorID
)

type SizeID uint

type Theme interface {
	Style(StyleID) lipgloss.Style
	Color(ColorID) lipgloss.Color
	Size(SizeID) int
}
