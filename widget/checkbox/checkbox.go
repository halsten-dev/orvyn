package checkbox

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/theme"
)

// Widget is a checkbox widget holding 2 state : Checked and Unchecked.
type Widget struct {
	orvyn.BaseWidget
	orvyn.BaseFocusable

	label string

	checked bool

	CheckKeybind key.Binding
}

// New creates and returns a new checkbox *Widget.
func New(label string) *Widget {
	w := new(Widget)

	w.BaseWidget = orvyn.NewBaseWidget()
	w.BaseFocusable = orvyn.NewBaseFocusable(w)

	w.checked = false
	w.label = label

	w.CheckKeybind = key.NewBinding(key.WithKeys(" "))

	w.OnBlur()

	return w
}

func (w *Widget) Update(msg tea.Msg) tea.Cmd {
	if m, ok := orvyn.GetKeyMsg(msg); ok {
		switch {
		case key.Matches(m, w.CheckKeybind):
			w.checked = !w.checked
		}
	}

	return nil
}

func (w *Widget) Resize(size orvyn.Size) {
	size.Height = 3

	w.BaseWidget.Resize(size)
}

func (w *Widget) Render() string {
	var checkbox string
	var label string
	var checked string

	style := w.GetStyle()

	checked = "   "

	if w.checked {
		checked = orvyn.GetTheme().Style(theme.TitleStyleID).Render(" X ")
	}

	checkbox = style.Render(checked)

	label = style.Width(w.GetContentSize().Width - 5).
		BorderStyle(lipgloss.HiddenBorder()).Render(w.label)

	return lipgloss.JoinHorizontal(lipgloss.Center, checkbox, label)
}

func (w *Widget) GetMinSize() orvyn.Size {
	return orvyn.NewSize(15, 3)
}

func (w *Widget) GetPreferredSize() orvyn.Size {
	return orvyn.NewSize(46, 3)
}

// IsChecked returns the current state of the checkbox.
func (w *Widget) IsChecked() bool {
	return w.checked
}

// SetChecked changes the current state of the checkbox.
func (w *Widget) SetChecked(checked bool) {
	w.checked = checked
}

// SetLabel changes the label associated to the checkbox.
func (w *Widget) SetLabel(label string) {
	w.label = label
}
