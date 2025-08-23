package theme

import (
	"github.com/charmbracelet/lipgloss"
)

type StyleName uint

const (
	TitleStyleName StyleName = iota
	NormalTextStyleName
	FocusedWidgetStyleName
	BlurredWidgetStyleName
	PaginatorActiveStyleName
	PaginatorInactiveStyleName
)

type ColorName uint

const (
	TitleFontColorName ColorName = iota
	NormalFontColorName
	DimFontColorName
	FocusedBorderColorName
	FocusedFontColorName
	BlurredBorderColorName
	BlurredFontColorName
)

type SizeName string

type Theme interface {
	Style(StyleName) lipgloss.Style
	Color(ColorName) lipgloss.Color
	Size(SizeName) int
}
