package orvyn

import tea "github.com/charmbracelet/bubbletea"

// Updatable interface gives ability to being updated.
type Updatable interface {
	Update(tea.Msg) tea.Cmd
}
