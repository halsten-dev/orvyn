package theme

import "github.com/charmbracelet/lipgloss"

type StyleName string

const (
	TitleStyleName         StyleName = "title_style_name"
	NormalTextName         StyleName = "normal_text_style_name"
	FocusedWidgetStyleName StyleName = "focused_widget_style_name"
	BlurredWidgetStyleName StyleName = "blurred_widget_style_name"
)

type ColorName string

const (
	TitleFontColorName     ColorName = "title_font_color_name"
	NormalFontColorName    ColorName = "normal_font_color_name"
	FocusedBorderColorName ColorName = "focused_border_color_name"
	FocusedFontColorName   ColorName = "focused_font_color_name"
	BlurredBorderColorName ColorName = "blurred_border_color_name"
	BlurredFontColorName   ColorName = "blurred_font_color_name"
)

type SizeName string

type Theme interface {
	Style(StyleName) lipgloss.Style
	Color(ColorName) lipgloss.Color
	Size(SizeName) int
}
