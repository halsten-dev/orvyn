package listdemo

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/layout"
	"github.com/halsten-dev/orvyn/widget/list"
)

type Screen struct {
	stringList   *list.Widget[string]
	stringValues []string

	layout *layout.CenterLayout
}

func New() *Screen {
	s := new(Screen)

	s.stringList = list.New(list.SimpleListItemConstructor)

	s.layout = layout.NewCenterLayout(s.stringList)

	return s
}

func (s *Screen) OnEnter(a any) tea.Cmd {
	s.stringValues = make([]string, 60)

	for i := 0; i < len(s.stringValues); i++ {
		mod := "a"

		if i%2 == 0 {
			mod = "b"
		}

		s.stringValues[i] = fmt.Sprintf("String %s %d", mod, i)
	}

	s.stringList.SetItems(s.stringValues)

	return nil
}

func (s *Screen) OnExit() any {
	return nil
}

func (s *Screen) Update(msg tea.Msg) tea.Cmd {
	s.stringList.Update(msg)

	return nil
}

func (s *Screen) Render() orvyn.Layout {
	return s.layout
}
