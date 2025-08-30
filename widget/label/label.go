package label

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/theme"
)

type Widget struct {
	orvyn.BaseWidget

	style lipgloss.Style

	value string
}

func New(value string) *Widget {
	w := new(Widget)

	w.BaseWidget = orvyn.NewBaseWidget()

	w.value = value

	return w
}

func (w *Widget) SetValue(value string) {
	w.value = value
}

func (w *Widget) Render() string {
	w.style = orvyn.GetTheme().Style(theme.LabelTextStyleID)

	size := w.GetSize()

	return w.style.
		Width(size.Width).
		Height(size.Height).
		Render(w.value)
}
