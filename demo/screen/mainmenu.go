package screen

import (
	"github.com/halsten-dev/orvyn/demo/widget"
	"github.com/halsten-dev/orvyn/widget/widgetlist"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/layout"
	"github.com/halsten-dev/orvyn/theme"
)

const IDMainMenu orvyn.ScreenID = "mainmenu"

type MainMenu struct {
	title *orvyn.SimpleRenderable

	actionList *widgetlist.Widget[widget.ActionListItemData]

	layout *layout.CenterLayout
}

func NewMainMenu() *MainMenu {
	m := new(MainMenu)

	m.title = orvyn.NewSimpleRenderable("Orvyn demo")
	m.title.Style = orvyn.GetTheme().Style(theme.TitleStyleID)

	menuItems := []widget.ActionListItemData{
		{
			Label:  "Input widget demo",
			Action: m.inputDemo,
		},
		{
			Label:  "WidgetList demo",
			Action: m.listDemo,
		},
		{
			Label:  "Progress demo",
			Action: m.progressDemo,
		},
		{
			Label:  "Quit",
			Action: m.quit,
		},
	}

	m.actionList = widgetlist.New(widget.ActionListItemConstructor)
	m.actionList.SetFilterable(false)
	m.actionList.SetMinSize(orvyn.NewSize(40, 4*len(menuItems)))
	m.actionList.SetPreferredSize(orvyn.NewSize(40, 4*len(menuItems)))

	m.actionList.SetItems(menuItems)

	m.layout = layout.NewCenterLayout(
		layout.NewVBoxLayout(10,
			m.title,
			orvyn.VGap,
			m.actionList,
		),
	)

	return m
}

func (m *MainMenu) inputDemo() tea.Cmd {
	return orvyn.SwitchScreen(InputWidgetDemoScreenID)
}

func (m *MainMenu) listDemo() tea.Cmd {
	return orvyn.SwitchScreen(ListDemoScreenID)
}

func (m *MainMenu) progressDemo() tea.Cmd {
	return orvyn.SwitchScreen(ProgressDemoScreenID)
}

func (m *MainMenu) quit() tea.Cmd {
	return tea.Quit
}

func (m *MainMenu) OnEnter(_ any) tea.Cmd {
	m.actionList.OnFocus()

	return nil
}

func (m *MainMenu) OnExit() any {
	return nil
}

func (m *MainMenu) Update(msg tea.Msg) tea.Cmd {
	if msg, ok := orvyn.GetKeyMsg(msg); ok {
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			item := m.actionList.GetSelectedItem()

			if item.Action != nil {
				return item.Action()
			}
		}
	}

	cmd := m.actionList.Update(msg)

	return cmd
}

func (m *MainMenu) Render() orvyn.Layout {
	return m.layout
}
