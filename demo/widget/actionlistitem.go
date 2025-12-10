package widget

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/widget/widgetlist"
)

type ActionListItemData struct {
	Label  string
	Action func() tea.Cmd
}

type ActionListItem struct {
	orvyn.BaseWidget
	orvyn.BaseFocusable

	data ActionListItemData
}

func ActionListItemConstructor(data ActionListItemData) widgetlist.ListItem[ActionListItemData] {
	a := new(ActionListItem)

	a.BaseWidget = orvyn.NewBaseWidget()
	a.BaseFocusable = orvyn.NewBaseFocusable(a)

	a.UpdateData(data)

	a.OnBlur()

	return a
}

// FilterValue implements list.ListItem.
func (a *ActionListItem) FilterValue() string {
	return ""
}

// GetData implements list.ListItem.
func (a *ActionListItem) GetData() ActionListItemData {
	return a.data
}

// OnEnterInput implements list.ListItem.
func (a *ActionListItem) OnEnterInput() {
}

// OnExitInput implements list.ListItem.
func (a *ActionListItem) OnExitInput() {
}

func (a *ActionListItem) Resize(size orvyn.Size) {
	size.Height = 3

	a.BaseWidget.Resize(size)
}

// Render implements list.ListItem.
func (a *ActionListItem) Render() string {
	contentSize := a.GetContentSize()

	return a.GetStyle().Width(contentSize.Width).
		Height(contentSize.Height).
		AlignHorizontal(lipgloss.Center).Render(a.data.Label)
}

// UpdateData implements list.ListItem.
func (a *ActionListItem) UpdateData(data ActionListItemData) {
	a.data = data
}
