package widgetlist

import (
	"github.com/halsten-dev/orvyn"
)

type SimpleListItem struct {
	orvyn.BaseWidget
	orvyn.BaseFocusable

	value string
}

func SimpleListItemConstructor(value string) ListItem[string] {
	sli := new(SimpleListItem)

	sli.BaseWidget = orvyn.NewBaseWidget()
	sli.BaseFocusable = orvyn.NewBaseFocusable(sli)

	sli.value = value

	sli.OnBlur()

	return sli
}

func (s *SimpleListItem) UpdateData(value string) {
	s.value = value
}

func (s *SimpleListItem) GetData() string {
	return s.value
}

func (s *SimpleListItem) Resize(size orvyn.Size) {
	size.Height = 3
	s.BaseWidget.Resize(size)
}

func (s *SimpleListItem) Render() string {
	size := s.GetContentSize()

	return s.GetStyle().
		Width(size.Width).
		Render(s.value)
}

func (s *SimpleListItem) OnEnterInput() {}

func (s *SimpleListItem) OnExitInput() {}

func (s *SimpleListItem) FilterValue() string {
	return s.value
}
