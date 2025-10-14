package label

import (
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/theme"
)

type Widget struct {
	orvyn.BaseWidget

	value string
}

func New(value string) *Widget {
	w := new(Widget)

	w.BaseWidget = orvyn.NewBaseWidget()
	w.BaseWidget.SetStyle(orvyn.GetTheme().Style(theme.LabelTextStyleID))
	w.value = value

	return w
}

func (w *Widget) SetValue(value string) {
	w.value = value
}

func (w *Widget) Render() string {
	size := w.BaseWidget.GetContentSize()

	return w.GetStyle().
		Width(size.Width).
		Height(size.Height).
		Render(w.value)
}
