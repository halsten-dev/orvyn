package dialog

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/layout"
	"github.com/halsten-dev/orvyn/theme"
	"github.com/halsten-dev/orvyn/widget/progressbar"
)

// Progress is a dialog for quick implementation of a progress dialog.
type Progress struct {
	progressBar     *progressbar.Widget
	srCancelKeybind *orvyn.SimpleRenderable

	maxSteps int
	steps    int
	percent  float64

	layout *layout.CenterLayout

	tickTag uint

	// cancelKeybind hold a *key.Binding that can be nil if no interruption is authorized.
	cancelKeybind *key.Binding

	// Interrupted flag will hold true if the progress was interrupted by the user.
	Interrupted bool
}

// NewProgress returns a new screen that represents a progress dialog.
// This screen needs to be used with orvyn.OpenDialog().
func NewProgress(title string) *Progress {
	p := &Progress{
		progressBar:     progressbar.New(title),
		srCancelKeybind: orvyn.NewSimpleRenderable(""),
	}

	p.layout = layout.NewCenterLayout(
		layout.NewMaxWidthVBoxLayout(10,
			p.progressBar,
			p.srCancelKeybind,
		),
	)

	p.cancelKeybind = nil
	p.srCancelKeybind.SetActive(false)
	p.Interrupted = false

	return p
}

func (p *Progress) OnEnter(i any) tea.Cmd {
	return orvyn.TickCmd(0, p.tickTag)
}

func (p *Progress) OnExit() any {
	return nil
}

func (p *Progress) Update(msg tea.Msg) tea.Cmd {
	cmd := p.progressBar.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if p.cancelKeybind != nil {
			if key.Matches(msg, *p.cancelKeybind) {
				p.Interrupted = true
				return orvyn.CloseDialog()
			}
		}
	case orvyn.TickMsg:
		if msg.Tag != p.tickTag {
			return nil
		}

		cmd := p.updateProgressBar()

		p.tickTag++
		return tea.Batch(cmd, orvyn.TickCmd(1, p.tickTag))
	}

	if p.percent >= 1 {
		return orvyn.CloseDialog()
	}

	return cmd
}

func (p *Progress) Render() orvyn.Layout {
	return p.layout
}

// Reset helps resetting the dialog to it's default state.
func (p *Progress) Reset() {
	p.Interrupted = false
	p.tickTag = 0
	p.steps = 0
}

// UpdateProgress should be used to update the underlying progressBar.
func (p *Progress) UpdateProgress(steps, maxSteps int) {
	var percent float64

	if maxSteps > 0 {
		percent = float64(100*steps/maxSteps) / 100
	} else {
		percent = 0
	}

	p.maxSteps = maxSteps
	p.steps = steps
	p.percent = percent
}

// SetCancelKeybind helps defining or removing a cancel keybind.
// The given keybind must have the Help initialized.
// For example :
//
// keybind := key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel"))
// p.dial.SetCancelKeybind(&keybind)
func (p *Progress) SetCancelKeybind(key *key.Binding) {
	p.cancelKeybind = key

	if key == nil {
		p.srCancelKeybind.SetValue("")
		p.srCancelKeybind.SetActive(false)
		return
	}

	keyText := key.Help().Key
	keyDesc := key.Help().Desc

	p.srCancelKeybind.SetValue(fmt.Sprintf("\n%s - %s",
		orvyn.GetTheme().Style(theme.HighlightTextStyleID).Render(keyText),
		keyDesc))
	p.srCancelKeybind.SetActive(true)
}

// SetBarColor helps changing the underlying progressBar color.
func (p *Progress) SetBarColor(color lipgloss.Color) {
	p.progressBar.SetColor(color)
}

func (p *Progress) updateProgressBar() tea.Cmd {
	p.progressBar.MaxValue = p.maxSteps
	p.progressBar.CurrentValue = p.steps

	return p.progressBar.SetPercent(p.percent)
}
