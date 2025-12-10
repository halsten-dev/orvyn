package screen

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/layout"
	"github.com/halsten-dev/orvyn/widget/checkbox"
	"github.com/halsten-dev/orvyn/widget/textarea"
	"github.com/halsten-dev/orvyn/widget/textinput"
)

type InputWidgetDemo struct {
	tiDemo *textinput.Widget
	taDemo *textarea.Widget
	cbDemo *checkbox.Widget

	focusManager *orvyn.FocusManager

	layout *layout.CenterLayout
}

func NewInputWidgetDemo() *InputWidgetDemo {
	s := new(InputWidgetDemo)

	s.tiDemo = textinput.New()

	s.taDemo = textarea.New()
	s.taDemo.SetMinSize(orvyn.NewSize(30, 1))
	s.taDemo.SetPreferredSize(orvyn.NewSize(30, 5))

	s.cbDemo = checkbox.New("Test")

	s.focusManager = orvyn.NewFocusManager()
	s.focusManager.Add(s.tiDemo)
	s.focusManager.Add(s.taDemo)
	s.focusManager.Add(s.cbDemo)

	s.layout = layout.NewCenterLayout(
		layout.NewMaxWidthVBoxFullLayout(orvyn.NewSize(10, 5), 1,
			s.tiDemo,
			s.taDemo,
			s.cbDemo,
		),
	)

	return s
}

func (s *InputWidgetDemo) OnEnter(a any) tea.Cmd {
	s.focusManager.FocusFirst()

	return nil
}

func (s *InputWidgetDemo) OnExit() any {
	return nil
}

func (s *InputWidgetDemo) Update(msg tea.Msg) tea.Cmd {
	if m, ok := orvyn.GetKeyMsg(msg); ok {
		switch {
		case key.Matches(m, key.NewBinding(key.WithKeys("esc"))):
			return orvyn.SwitchToPreviousScreen()
		}
	}

	cmd := s.focusManager.Update(msg)

	return cmd
}

func (s *InputWidgetDemo) Render() orvyn.Layout {
	return s.layout
}
