package theme

import (
	"github.com/charmbracelet/lipgloss"
)

type StyleID uint

const (
	TitleStyleID StyleID = iota
	NormalTextStyleID
	FocusedWidgetStyleID
	BlurredWidgetStyleID
	PaginatorActiveStyleID
	PaginatorInactiveStyleID
	TextInputTextStyleID
	TextInputCursorStyleID
	TextInputCursorTextStyleID
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
)

type SizeID string

type Theme interface {
	Style(StyleID) lipgloss.Style
	Color(ColorID) lipgloss.Color
	Size(SizeID) int
}
