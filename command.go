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

// TickCmd should be returned in update function to ensure the tick continues.
// Be sure to check and increment the tick tag when responding to the message.
//
//	func (s *Screen) Update(msg tea.Msg) tea.Cmd {
//		switch msg := msg.(type) {
//		case orvyn.TickMsg:
//			if msg.Tag != s.tickTag {
//				return nil
//			}
//
//			s.updateData()
//
//			s.tickTag++
//			return orvyn.TickCmd(tick, s.tickTag)
//		}
//	}
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

func dialogExitCmd(id ScreenID, param any) tea.Cmd {
	return func() tea.Msg {
		return DialogExitMsg{
			DialogID: id,
			Param:    param,
		}
	}
}
