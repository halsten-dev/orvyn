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

func (s *SimpleRenderable) GetMinSize() Size {
	vFrame, hFrame := s.Style.GetFrameSize()

	if vFrame+hFrame == 0 {
		s.BaseRenderable.GetMinSize()
	}

	return GetRenderSize(s.Style, "min")
}

func (s *SimpleRenderable) GetPreferredSize() Size {
	vFrame, hFrame := s.Style.GetFrameSize()

	if vFrame+hFrame == 0 {
		s.BaseRenderable.GetPreferredSize()
	}

	return GetRenderSize(s.Style, "pref")
}
