package textinput

import (
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/theme"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Widget struct {
	orvyn.BaseWidget
	orvyn.BaseFocusable

	textinput.Model
}

func New() *Widget {
	w := new(Widget)

	t := orvyn.GetTheme()

	w.BaseWidget = orvyn.NewBaseWidget()
	w.BaseFocusable = orvyn.NewBaseFocusable(w)

	w.Model = textinput.New()
	w.Prompt = ""
	w.TextStyle = t.Style(theme.NormalTextStyleID)
	w.Cursor.Style = t.Style(theme.NormalTextStyleID)
	w.Cursor.TextStyle = t.Style(theme.NormalTextStyleID)

	w.OnBlur()

	return w
}

func (w *Widget) Init() tea.Cmd {
	w.Model.SetValue("")
	return textinput.Blink
}

func (w *Widget) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	w.Model, cmd = w.Model.Update(msg)

	return cmd
}

func (w *Widget) OnFocus() {
	w.BaseFocusable.OnFocus()
	w.Model.Focus()
}

func (w *Widget) OnBlur() {
	w.BaseFocusable.OnBlur()
	w.Model.Blur()
}

func (w *Widget) Render() string {
	return w.GetStyle().Render(w.Model.View())
}

func (w *Widget) Resize(size orvyn.Size) {
	style := w.GetStyle()
	size.Height = 1 + style.GetVerticalFrameSize()

	w.BaseWidget.Resize(size)

	contentSize := w.GetContentSize()
	// Take borders into account
	w.Model.Width = contentSize.Width - style.GetHorizontalFrameSize()
	w.Model.Width -= max(0, len(w.Model.Prompt))
	w.Model.Width = max(2, w.Model.Width)

	// For the Bubbles textinput to process the update
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
	return orvyn.NewSize(26, 3)
}

func (w *Widget) GetPreferredSize() orvyn.Size {
	return orvyn.NewSize(46, 3)
}

func (w *Widget) OnEnterInput() {}

func (w *Widget) OnExitInput() {}
