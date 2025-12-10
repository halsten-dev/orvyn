package orvyn

import (
	tea "github.com/charmbracelet/bubbletea"
)

// ScreenID type represents an ID for a screen.
type ScreenID string

// Screen interface defines behaviour of an Orvyn Screen.
type Screen interface {
	// OnEnter is called when the screen is entered. Can take as parameter a struct from the previous screen.
	OnEnter(any) tea.Cmd

	// OnExit is called when the screen is being exited. Can return a struct that will be passed to the next screen.
	OnExit() any

	// Updatable so the Screen can be updated.
	Updatable

	// Render returns the view string of the whole screen.
	Render() Layout
}
