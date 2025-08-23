package orvyn

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

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
