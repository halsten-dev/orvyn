package dialog

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/layout"
	"github.com/halsten-dev/orvyn/theme"
)

type Option struct {
	Keybind key.Binding
	Text    string
	Value   uint
}

type Config struct {
	Message string
	Options []Option
}

type Popup struct {
	config Config

	content *orvyn.SimpleRenderable
	options *orvyn.SimpleRenderable

	layout *layout.CenterLayout

	value uint
}

// NewPopup returns a new screen based on the given Config.
// This screen needs to be used with orvyn.OpenDialog().
func NewPopup(config Config) *Popup {
	var b strings.Builder

	s := new(Popup)

	t := orvyn.GetTheme()
	ns := t.Style(theme.NormalTextStyleID)
	ds := t.Style(theme.DimTextStyleID)
	nds := t.Style(theme.NeutralDimTextStyleID)

	s.config = config

	b.WriteString(config.Message)
	b.WriteString("\n\n")

	s.content = orvyn.NewSimpleRenderable(b.String())
	s.content.Style = ns.AlignHorizontal(lipgloss.Center)
	s.content.SizeConstraint = true

	b.Reset()

	for i, o := range config.Options {
		if i > 0 {
			b.WriteString(nds.Render(fmt.Sprintf(" %c ", '•')))
		}

		fmt.Fprintf(&b, "%s %s",
			ns.Render(o.Keybind.Help().Key),
			ds.Render(o.Text))
	}

	s.options = orvyn.NewSimpleRenderable(b.String())

	s.layout = layout.NewCenterLayout(
		layout.NewVBoxLayout(10,
			s.content,
			s.options,
		),
	)

	return s
}

func (s *Popup) OnEnter(i any) tea.Cmd {
	return nil
}

func (s *Popup) OnExit() any {
	return s.value
}

func (s *Popup) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		for _, o := range s.config.Options {
			if key.Matches(msg, o.Keybind) {
				s.value = o.Value
				return orvyn.CloseDialog()
			}
		}
	}

	return nil
}

func (s *Popup) Render() orvyn.Layout {
	return s.layout
}
