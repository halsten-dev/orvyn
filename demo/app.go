package main

import (
	"github.com/halsten-dev/orvyn"

	tea "github.com/charmbracelet/bubbletea"
)

// App is the main model to run the Orvyn application
type App struct{}

func (a App) Init() tea.Cmd {
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := orvyn.Update(msg)

	return a, cmd
}

func (a App) View() string {
	return orvyn.Render()
}
