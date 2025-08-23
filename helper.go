package orvyn

import (
	"github.com/charmbracelet/lipgloss"
)

func GetRenderSize(style lipgloss.Style, value string) Size {
	width, height := lipgloss.Size(style.Render(value))

	return NewSize(width, height)
}
