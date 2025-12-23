package screen

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/dialog"
	"github.com/halsten-dev/orvyn/layout"
)

type ProgressDemo struct {
	dial *dialog.Progress

	layout *layout.CenterLayout

	srInstruction *orvyn.SimpleRenderable
	srStatus      *orvyn.SimpleRenderable
}

func NewProgressDemo() *ProgressDemo {
	p := &ProgressDemo{
		dial:          dialog.NewProgress("On going"),
		srInstruction: orvyn.NewSimpleRenderable("Press <Space> to launch progress"),
		srStatus:      orvyn.NewSimpleRenderable("Progress finished !"),
	}

	p.srStatus.SetActive(false)

	keybind := key.NewBinding(key.WithKeys("esc"))
	p.dial.CancelKeybind = &keybind

	p.layout = layout.NewCenterLayout(
		layout.NewMaxWidthVBoxLayout(0,
			p.srInstruction,
			p.srStatus,
		),
	)

	return p
}

func (p *ProgressDemo) OnEnter(i any) tea.Cmd {
	return nil
}

func (p *ProgressDemo) OnExit() any {
	return nil
}

func (p *ProgressDemo) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys(" "))):
			return p.launchProgress()
		}
	case orvyn.DialogExitMsg:
		switch msg.DialogID {
		case "progressBar":
			if p.dial.Interrupted {
				p.srStatus.SetValue("Interrupted !")
			} else {
				p.srStatus.SetValue("Progress finished !")
			}

			p.srStatus.SetActive(true)
		}
	}

	return nil
}

func (p *ProgressDemo) Render() orvyn.Layout {
	return p.layout
}

func (p *ProgressDemo) launchProgress() tea.Cmd {

	// Loop through every keys
	p.dial.Reset()

	go func(dial *dialog.Progress) {
		count := 0
		maxSteps := 100

		dial.UpdateProgress(count, maxSteps)

		for range maxSteps {
			if dial.Interrupted {
				return
			}

			count++
			dial.UpdateProgress(count, maxSteps)

			time.Sleep(300 * time.Millisecond)
		}
	}(p.dial)

	return orvyn.OpenDialog("progressBar", p.dial, nil)
}
