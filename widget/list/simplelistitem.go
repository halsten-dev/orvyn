package list

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/theme"
)

type SimpleListItem struct {
	orvyn.BaseWidget
	orvyn.BaseFocusable
	style lipgloss.Style

	value string
}

func SimpleListItemConstructor(value string) IListItem {
	sli := new(SimpleListItem)

	sli.BaseWidget = orvyn.NewBaseWidget()

	sli.value = value

	sli.OnBlur()

	return sli
}

func (s *SimpleListItem) Resize(size orvyn.Size) {
	size.Width -= s.style.GetHorizontalFrameSize()
	size.Height = lipgloss.Height(s.style.Render(s.value))

	s.BaseWidget.Resize(size)
}

func (s *SimpleListItem) Render() string {
	size := s.GetSize()

	return s.style.
		Width(size.Width).
		Render(s.value)
}

func (s *SimpleListItem) OnFocus() {
	s.style = orvyn.GetTheme().Style(theme.FocusedWidgetStyleName)
}

func (s *SimpleListItem) OnBlur() {
	s.style = orvyn.GetTheme().Style(theme.BlurredWidgetStyleName)
}

func (s *SimpleListItem) OnEnterInput() {}

func (s *SimpleListItem) OnExitInput() {}
