package listdemo

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/layout"
	"github.com/halsten-dev/orvyn/widget/list"
)

type Screen struct {
	stringList   *list.Widget[string]
	stringValues []string

	elementIndex int

	layout *layout.CenterLayout
}

func New() *Screen {
	s := new(Screen)

	s.stringList = list.New(list.SimpleListItemConstructor)
	s.stringList.AutoFocusNewItem = false

	s.elementIndex = 0

	s.layout = layout.NewCenterLayout(s.stringList)

	return s
}

func (s *Screen) OnEnter(a any) tea.Cmd {
	s.stringValues = make([]string, 0)

	s.elementIndex = 0

	s.stringList.SetItems(s.stringValues)

	return nil
}

func (s *Screen) OnExit() any {
	return nil
}

func (s *Screen) Update(msg tea.Msg) tea.Cmd {
	if m, ok := orvyn.GetKeyMsg(msg); ok {
		switch {
		case key.Matches(m, key.NewBinding(key.WithKeys("n"))):
			if s.stringList.FilterState() != list.Filtering {
				s.stringList.AppendItem(fmt.Sprintf("Test %d", s.elementIndex))
				s.elementIndex++
			}

		case key.Matches(m, key.NewBinding(key.WithKeys("i"))):
			if s.stringList.FilterState() != list.Filtering {
				s.stringList.InsertItem(s.stringList.GetGlobalIndex(), fmt.Sprintf("Test Insert %d", s.elementIndex))
				s.elementIndex++
			}

		case key.Matches(m, key.NewBinding(key.WithKeys("shift+up"))):
			currentIndex := s.stringList.GetGlobalIndex()
			if currentIndex > 0 {
				s.stringList.MoveItem(currentIndex, currentIndex-1)
			}

		case key.Matches(m, key.NewBinding(key.WithKeys("shift+down"))):
			currentIndex := s.stringList.GetGlobalIndex()
			if currentIndex < s.stringList.Length()-1 {
				s.stringList.MoveItem(currentIndex, currentIndex+1)
			}
		}
	}

	cmd := s.stringList.Update(msg)

	return cmd
}

func (s *Screen) Render() orvyn.Layout {
	return s.layout
}
