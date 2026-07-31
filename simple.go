package orvyn

import (
	"github.com/charmbracelet/lipgloss"
)

// SimpleRenderable represents a very basic and reusable renderable.
type SimpleRenderable struct {
	// BaseRenderable composing
	BaseRenderable

	// Style is the style that will be used to render the value of the SimpleRenderable.
	Style lipgloss.Style

	// SizeConstraint define if the style width and height should be taken into account when rendering the SimpleRenderable.
	SizeConstraint bool

	// value is the value of the SimpleRenderable.
	value string
}

// VGap is a SimpleRenderable representing a new line. Useful for layout building.
var VGap = NewSimpleRenderable("\n")

// NewSimpleRenderable creates a new SimpleRenderable and returns it.
func NewSimpleRenderable(value string) *SimpleRenderable {
	s := new(SimpleRenderable)

	s.BaseRenderable = NewBaseRenderable()
	s.value = value
	s.Style = lipgloss.NewStyle()
	s.SizeConstraint = false

	return s
}

// SetValue changes the current value with the given one.
func (s *SimpleRenderable) SetValue(value string) {
	s.value = value
}

func (s *SimpleRenderable) Render() string {
	if !s.SizeConstraint {
		return s.Style.Render(s.value)
	}

	size := s.GetSize()

	size.Width -= s.Style.GetHorizontalFrameSize()
	size.Height -= s.Style.GetVerticalFrameSize()

	return s.Style.Width(size.Width).
		Height(size.Height).Render(s.value)
}

// GetMinSize returns the size the value actually renders at, unless a size was
// set explicitly with SetMinSize.
//
// It used to measure the literal string "min", so every SimpleRenderable claimed
// to be one line tall whatever it held. Layouts budgeted one row for a multi-line
// value and the overflow ate their margin.
func (s *SimpleRenderable) GetMinSize() Size {
	if size := s.BaseRenderable.GetMinSize(); size != NewSize(1, 1) {
		return size
	}

	return GetRenderSize(s.Style, s.value)
}

// GetPreferredSize returns the size the value actually renders at, unless a size
// was set explicitly with SetPreferredSize. See GetMinSize.
func (s *SimpleRenderable) GetPreferredSize() Size {
	if size := s.BaseRenderable.GetPreferredSize(); size != NewSize(1, 1) {
		return size
	}

	return GetRenderSize(s.Style, s.value)
}
