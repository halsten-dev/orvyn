package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/demo/screen"
)

func main() {
	// Orvyn
	orvyn.Init()

	orvyn.RegisterScreen(screen.MainMenuScreenID, screen.NewMainMenu())
	orvyn.RegisterScreen(screen.ListDemoScreenID, screen.NewListDemo())
	orvyn.RegisterScreen(screen.InputWidgetDemoScreenID, screen.NewInputWidgetDemo())
	orvyn.RegisterScreen(screen.ProgressDemoScreenID, screen.NewProgressDemo())

	p := tea.NewProgram(&App{}, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
