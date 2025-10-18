package orvyn

import (
	"github.com/charmbracelet/lipgloss"
)

type SimpleRenderable struct {
	BaseRenderable

	Style          lipgloss.Style
	SizeConstraint bool

	value string
}

var VGap = NewSimpleRenderable("\n")

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
