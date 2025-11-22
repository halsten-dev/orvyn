package orvyn

import (
	"github.com/charmbracelet/lipgloss"
)

// GetRenderSize returns the Size of a value drawn with a given style.
func GetRenderSize(style lipgloss.Style, value string) Size {
	width, height := lipgloss.Size(style.Render(value))

	return NewSize(width, height)
}
