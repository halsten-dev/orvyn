package main

import (
	"fmt"

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

	size1, size2 := orvyn.DivideSizeFull(93)

	fmt.Printf("Size 1 = %d, size 2 = %d", size1, size2)

	// p := tea.NewProgram(&App{}, tea.WithAltScreen())
	//
	// if _, err := p.Run(); err != nil {
	// 	log.Fatal(err)
	// }
}
