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

	MinHeight       int
	PreferredHeight int
}

func New() *Widget {
	w := new(Widget)

	w.BaseWidget = orvyn.NewBaseWidget()

	w.Model = textarea.New()
	w.Model.Prompt = ""
	w.Model.SetWidth(10)

	w.OnBlur()

	w.MinHeight = 1
	w.PreferredHeight = 5

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
	w.updateStyle()
	w.Model.Focus()
}

func (w *Widget) OnBlur() {
	w.updateStyle()
	w.Model.Blur()
}

func (w *Widget) Render() string {
	return w.Model.View()
}

func (w *Widget) Resize(size orvyn.Size) {
	w.BaseWidget.Resize(size)

	contentSize := w.GetContentSize()

	w.Model.SetWidth(size.Width)
	w.Model.SetHeight(contentSize.Height)

	focused := w.Model.Focused()
	if !focused {
		w.Model.Focus()
	}

	w.Model, _ = w.Model.Update(nil)

	if !focused {
		w.Model.Blur()
	}
}

func (w *Widget) GetMinSize() orvyn.Size {
	return orvyn.NewSize(20, w.MinHeight)
}

func (w *Widget) GetPreferredSize() orvyn.Size {
	return orvyn.NewSize(50, w.PreferredHeight)
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
