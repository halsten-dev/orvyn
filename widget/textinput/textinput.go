package textinput

import (
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/theme"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Widget struct {
	orvyn.BaseWidget
	orvyn.BaseFocusable

	style lipgloss.Style

	textinput.Model
}

func New() *Widget {
	w := new(Widget)

	w.BaseWidget = orvyn.NewBaseWidget()

	w.Model = textinput.New()
	w.Prompt = ""

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
	w.Model.Focus()
	w.updateStyle(true)
}

func (w *Widget) OnBlur() {
	w.Model.Blur()
	w.updateStyle(false)
}

func (w *Widget) Render() string {
	return w.style.Render(w.Model.View())
}

func (w *Widget) Resize(size orvyn.Size) {
	size.Height = 1 + w.style.GetVerticalFrameSize()

	w.BaseWidget.Resize(size)

	// Take borders into account
	w.Model.Width = size.Width - w.style.GetHorizontalFrameSize()
	w.Model.Width -= max(1, len(w.Model.Prompt))
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

func (w *Widget) updateStyle(focused bool) {
	t := orvyn.GetTheme()

	if focused {
		w.style = t.Style(theme.FocusedWidgetStyleID)
	} else {
		w.style = t.Style(theme.BlurredWidgetStyleID)
	}

	w.TextStyle = t.Style(theme.NormalTextStyleID)
	w.Cursor.Style = t.Style(theme.NormalTextStyleID)
	w.Cursor.TextStyle = t.Style(theme.NormalTextStyleID)
}
