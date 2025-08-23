package orvyn

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ScreenID string

// Screen defines behaviour of an Orvyn screen.
type Screen interface {
	// OnEnter is called when the screen is entered. Can take as parameter the struct from the previous screen.
	OnEnter(any) tea.Cmd

	// OnExit is called when the screen is being exited. Can return a struct that will be passed to the next screen.
	OnExit() any

	// Updatable Screen can be updated.
	Updatable

	// Render returns the view string of the whole screen
	Render() Layout
}
