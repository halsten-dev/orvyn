package dialog

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/layout"
	"github.com/halsten-dev/orvyn/widget/progressbar"
)

type Progress struct {
	progressBar *progressbar.Widget

	maxSteps int
	steps    int
	percent  float64

	layout *layout.CenterLayout

	tickTag uint

	CancelKeybind *key.Binding
	Interrupted   bool
}

func NewProgress(title string) *Progress {
	p := &Progress{
		progressBar: progressbar.New(title),
	}

	p.layout = layout.NewCenterLayout(
		p.progressBar,
	)

	p.CancelKeybind = nil
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
		if p.CancelKeybind != nil {
			if key.Matches(msg, *p.CancelKeybind) {
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

func (p *Progress) Reset() {
	p.Interrupted = false
	p.tickTag = 0
	p.steps = 0
}

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

func (p *Progress) updateProgressBar() tea.Cmd {
	p.progressBar.MaxValue = p.maxSteps
	p.progressBar.CurrentValue = p.steps

	return p.progressBar.SetPercent(p.percent)
}
