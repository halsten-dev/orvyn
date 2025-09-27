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

type Screen struct {
	config Config

	content *orvyn.SimpleRenderable
	options *orvyn.SimpleRenderable

	layout *layout.CenterLayout

	value uint
}

// New returns a new screen based on the given Config.
// This screen will need to be used with orvyn.OpenDialog().
func NewPopup(config Config) *Screen {
	var b strings.Builder

	s := new(Screen)

	t := orvyn.GetTheme()
	ns := t.Style(theme.NormalTextStyleID)
	ds := t.Style(theme.DimTextStyleID)

	s.config = config

	b.WriteString(config.Message)
	b.WriteString("\n\n")

	s.content = orvyn.NewSimpleRenderable(b.String())
	s.content.Style = ns.AlignHorizontal(lipgloss.Center)
	s.content.SizeConstraint = true

	b.Reset()

	for i, o := range config.Options {
		if i > 0 {
			b.WriteString(fmt.Sprintf(" %c ", '•'))
		}

		b.WriteString(fmt.Sprintf("%s %s",
			ns.Render(o.Keybind.Help().Key),
			ds.Render(o.Text)))
	}

	s.options = orvyn.NewSimpleRenderable(b.String())

	s.layout = layout.NewCenterLayout(
		layout.NewVBoxLayout(10,
			[]orvyn.Renderable{
				s.content,
				s.options,
			},
		),
	)

	return s
}

func (s *Screen) OnEnter(i any) tea.Cmd {
	return nil
}

func (s *Screen) OnExit() any {
	return s.value
}

func (s *Screen) Update(msg tea.Msg) tea.Cmd {
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

func (s *Screen) Render() orvyn.Layout {
	return s.layout
}
