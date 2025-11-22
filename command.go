package orvyn

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TickMsg simplify and assure the right tick msg management with a unique Tag.
type TickMsg struct {
	Time time.Time
	Tag  uint
}

func TickCmd(seconds time.Duration, tag uint) tea.Cmd {
	return tea.Tick(seconds*time.Second, func(t time.Time) tea.Msg {
		return TickMsg{
			Time: t,
			Tag:  tag,
		}
	})
}

// DialogExitMsg is the message sent when an Orvyn dialog is exited.
type DialogExitMsg struct {
	DialogID ScreenID
	Param    any
}

func DialogExitCmd(id ScreenID, param any) tea.Cmd {
	return func() tea.Msg {
		return DialogExitMsg{
			DialogID: id,
			Param:    param,
		}
	}
}
