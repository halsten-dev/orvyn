package inputwidgetdemo

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/layout"
	"github.com/halsten-dev/orvyn/widget/checkbox"
	"github.com/halsten-dev/orvyn/widget/textarea"
	"github.com/halsten-dev/orvyn/widget/textinput"
)

type Screen struct {
	tiDemo *textinput.Widget
	taDemo *textarea.Widget
	cbDemo *checkbox.Widget

	focusManager *orvyn.FocusManager

	layout *layout.CenterLayout
}

func New() *Screen {
	s := new(Screen)

	s.tiDemo = textinput.New()
	s.taDemo = textarea.New()
	s.cbDemo = checkbox.New("Test")

	s.focusManager = orvyn.NewFocusManager()
	s.focusManager.Add(s.tiDemo)
	s.focusManager.Add(s.taDemo)
	s.focusManager.Add(s.cbDemo)

	s.layout = layout.NewCenterLayout(
		layout.NewDefinedWidthVerticalLayout(30, 120, 10,
			[]orvyn.Renderable{
				s.tiDemo,
				s.taDemo,
				s.cbDemo,
			},
		),
	)

	return s
}

func (s *Screen) OnEnter(a any) tea.Cmd {
	s.focusManager.FocusFirst()

	return nil
}

func (s *Screen) OnExit() any {
	return nil
}

func (s *Screen) Update(msg tea.Msg) tea.Cmd {
	cmd := s.focusManager.Update(msg)

	return cmd
}

func (s *Screen) Render() orvyn.Layout {
	return s.layout
}
