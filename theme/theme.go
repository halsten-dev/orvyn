package theme

import (
	"github.com/charmbracelet/lipgloss"
)

type StyleID uint

const (
	TitleStyleID StyleID = iota
	NormalTextStyleID
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

type SizeID string

type Theme interface {
	Style(StyleID) lipgloss.Style
	Color(ColorID) lipgloss.Color
	Size(SizeID) int
}
