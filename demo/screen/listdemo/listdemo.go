package listdemo

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/layout"
	"github.com/halsten-dev/orvyn/widget/widgetlist"
)

type Screen struct {
	stringList   *widgetlist.Widget[string]
	stringValues []string

	elementIndex int

	layout *layout.CenterLayout
}

func New() *Screen {
	s := new(Screen)

	s.stringList = widgetlist.New(widgetlist.SimpleListItemConstructor)
	s.stringList.AutoFocusNewItem = false
	// s.stringList.Filter = func(items *[]widgetlist.ListItem[string], s string) widgetlist.FilteredItems {
	// 	var filteredItems widgetlist.FilteredItems
	//
	// 	for i, item := range *items {
	// 		if !strings.Contains(strings.ToLower(item.FilterValue()), strings.ToLower(s)) {
	// 			filteredItems = append(filteredItems, widgetlist.FilteredItem{
	// 				Index: i,
	// 			})
	// 		}
	// 	}
	//
	// 	return filteredItems
	// }

	s.elementIndex = 0

	s.layout = layout.NewCenterLayout(s.stringList)

	return s
}

func (s *Screen) OnEnter(a any) tea.Cmd {
	s.stringValues = []string{
		"Pêche au poisson",
		"Couper du bois",
		"Cueillir des champignons vénéneux dans les bois",
		"Courir tout nu dans le pré",
		"Se fritter avec des monstres gigantesque",
		"S'habiller avec des habits grotesque",
		"Farmer le blé dans le pré",
		"Dancer tout nu avec un epis de blé coincé entre les parties charnues",
		"S'en aller faire des courses",
	}

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
			if s.stringList.FilterState() != widgetlist.Filtering {
				s.stringList.AppendItem(fmt.Sprintf("Test %d", s.elementIndex))
				s.elementIndex++
			}

		case key.Matches(m, key.NewBinding(key.WithKeys("i"))):
			if s.stringList.FilterState() != widgetlist.Filtering {
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
