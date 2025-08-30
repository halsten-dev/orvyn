package statusmessage

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/theme"
)

type messageType int

const (
	ErrorMessage messageType = iota
	SuccessMessage
	WarningMessage
	InformationMessage
	NeutralMessage
)

type Widget struct {
	orvyn.BaseWidget

	message      string
	messageType  messageType
	messageStyle lipgloss.Style
}

func New() *Widget {
	w := new(Widget)

	w.BaseWidget = orvyn.NewBaseWidget()

	return w
}

func (w *Widget) Init() tea.Cmd {
	w.Reset()

	return nil
}

func (w *Widget) Render() string {
	size := w.GetSize()

	s := ""

	if w.message != "" {
		s = w.messageStyle.
			Width(size.Width).
			Render(w.message)
	}

	return s
}

func (w *Widget) GetMinSize() orvyn.Size {
	return orvyn.GetRenderSize(w.messageStyle, w.message)
}

func (w *Widget) GetPreferredSize() orvyn.Size {
	return w.GetMinSize()
}

func (w *Widget) SetMessage(msg string, msgType messageType) {
	w.message = msg
	w.messageType = msgType
	w.updateStyle()
}

func (w *Widget) SetError(err error) {
	w.message = err.Error()
	w.messageType = ErrorMessage
	w.updateStyle()
}

func (w *Widget) Reset() {
	w.message = ""
	w.messageType = NeutralMessage
	w.updateStyle()
}

func (w *Widget) updateStyle() {
	switch w.messageType {
	case ErrorMessage:
		w.messageStyle = orvyn.GetTheme().Style(theme.StatusErrorTextStyleID)
	case SuccessMessage:
		w.messageStyle = orvyn.GetTheme().Style(theme.StatusSuccessTextStyleID)
	case WarningMessage:
		w.messageStyle = orvyn.GetTheme().Style(theme.StatusWarningTextStyleID)
	case InformationMessage:
		w.messageStyle = orvyn.GetTheme().Style(theme.StatusInformationTextStyleID)
	case NeutralMessage:
		w.messageStyle = orvyn.GetTheme().Style(theme.StatusNeutralTextStyleID)
	}
}
