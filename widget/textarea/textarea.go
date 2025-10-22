package textarea

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/theme"
)

type Widget struct {
	orvyn.BaseWidget
	orvyn.BaseFocusable

	textarea.Model
}

func New() *Widget {
	w := new(Widget)

	w.BaseWidget = orvyn.NewBaseWidget()
	w.BaseFocusable = orvyn.NewBaseFocusable(w)

	w.Model = textarea.New()
	w.Model.Prompt = ""
	w.Model.SetWidth(10)

	w.OnBlur()

	return w
}

func (w *Widget) Init() tea.Cmd {
	w.Model.SetValue("")
	return textarea.Blink
}

func (w *Widget) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	w.Model, cmd = w.Model.Update(msg)

	return cmd
}

func (w *Widget) OnFocus() {
	w.BaseFocusable.OnFocus()
	w.updateStyle()
	w.Model.Focus()
}

func (w *Widget) OnBlur() {
	w.BaseFocusable.OnBlur()
	w.updateStyle()
	w.Model.Blur()
}

func (w *Widget) Render() string {
	return w.Model.View()
}

func (w *Widget) Resize(size orvyn.Size) {
	w.BaseWidget.Resize(size)

	w.Model.SetWidth(size.Width)
	w.Model.SetHeight(size.Height)

	focused := w.Model.Focused()
	if !focused {
		w.Model.Focus()
	}

	w.Model, _ = w.Model.Update(nil)

	if !focused {
		w.Model.Blur()
	}
}

func (w *Widget) OnEnterInput() {}

func (w *Widget) OnExitInput() {}

func (w *Widget) updateStyle() {
	t := orvyn.GetTheme()

	w.BlurredStyle.Text = t.Style(theme.NormalTextStyleID)
	w.BlurredStyle.Base = t.Style(theme.BlurredWidgetStyleID)
	w.BlurredStyle.CursorLine = t.Style(theme.NormalTextStyleID)
	w.FocusedStyle.Text = t.Style(theme.NormalTextStyleID)
	w.FocusedStyle.Base = t.Style(theme.FocusedWidgetStyleID)
	w.FocusedStyle.CursorLine = t.Style(theme.NormalTextStyleID)
	w.Cursor.TextStyle = t.Style(theme.NormalTextStyleID)
	w.Cursor.Style = t.Style(theme.NormalTextStyleID)
}
