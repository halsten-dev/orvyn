package main

import (
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/demo/screen"

	tea "github.com/charmbracelet/bubbletea"
)

// App is the main model to run the Orvyn application
type App struct{}

func (a App) Init() tea.Cmd {
	return orvyn.SwitchScreen(screen.ListDemoScreenID)
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := orvyn.Update(msg)

	return a, cmd
}

func (a App) View() string {
	return orvyn.Render()
}
