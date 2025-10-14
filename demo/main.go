package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/demo/screen"
	"github.com/halsten-dev/orvyn/demo/screen/inputwidgetdemo"
	"github.com/halsten-dev/orvyn/demo/screen/listdemo"
)

func main() {
	// Orvyn
	orvyn.Init()

	// orvyn.RegisterScreen(screen.IDProjectLoading, projectloading.New())
	orvyn.RegisterScreen(screen.ListDemoScreenID, listdemo.New())
	orvyn.RegisterScreen(screen.InputWidgetDemoScreenID, inputwidgetdemo.New())

	p := tea.NewProgram(&App{}, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
