package orvyn

import (
	"slices"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// FocusManager can be instantiated when needed on a Screen to manage focus
// on multiple widgets. Manages focus and input mode of registred widgets.
type FocusManager struct {
	// NextFocusKeybind holds the key.Binding to loop through the Focusable Widgets.
	// Tab by default.
	NextFocusKeybind key.Binding

	// PreviousFocusKeybind holds the key.Binding to loop through the Focusable Widgets.
	// Shift+Tab by default.
	PreviousFocusKeybind key.Binding

	// ManageFocusNextPrevKeybind grants the focusManager to react to the NextFocusKeybind and PreviousFocusKeybind.
	// True by default.
	ManageFocusNextPrevKeybind bool

	widgets     []Focusable
	tabIndex    int
	isInputting bool
}

// NewFocusManager creates and return a new *FocusManager.
func NewFocusManager() *FocusManager {
	f := new(FocusManager)

	f.widgets = make([]Focusable, 0)
	f.tabIndex = 0
	f.isInputting = false

	f.NextFocusKeybind = key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next focus"),
	)
	f.PreviousFocusKeybind = key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "previous focus"),
	)

	f.ManageFocusNextPrevKeybind = true

	return f
}

// Add appends the given Focusable Widget to the manager.
// Order of append defines the focus order.
func (f *FocusManager) Add(widget Focusable) {
	if slices.Contains(f.widgets, widget) {
		return
	}

	f.widgets = append(f.widgets, widget)
}

// Insert adds a Focusable Widget at the given index. Change the focus order.
func (f *FocusManager) Insert(index int, widget Focusable) {
	if index < 0 || index >= len(f.widgets) {
		return
	}

	if slices.Contains(f.widgets, widget) {
		return
	}

	f.widgets = append(f.widgets[:index+1], f.widgets[index:]...)
	f.widgets[index] = widget

	if index <= f.tabIndex {
		f.tabIndex = f.getNextIndex()
	}
}

// Update allows to change the widget at the given index, whitout changing the focus order.
func (f *FocusManager) UpdateWidget(index int, widget Focusable) {
	if index < 0 || index >= len(f.widgets) {
		return
	}

	f.widgets[index] = widget
}

// Remove the widget from the manager at the given index.
// Automatically focus the previous widget, if the given index is the currently focused widget.
func (f *FocusManager) Remove(index int) {
	if index < 0 || index >= len(f.widgets) {
		return
	}

	f.widgets = append(f.widgets[:index], f.widgets[index+1:]...)

	if f.tabIndex == index {
		f.tabIndex = f.getPreviousIndex()
	}
}

// RemoveWidget removes the widget from the manager base on the given Focusable.
// Need to be a pointer for this function to work as expected.
func (f *FocusManager) RemoveWidget(widget Focusable) {
	for i, w := range f.widgets {
		if w == widget {
			f.Remove(i)
			return
		}
	}
}

// SetWidgets replaces the manager widget list with the one given.
// Given list order defines the focus order.
func (f *FocusManager) SetWidgets(widgets []Focusable) {
	f.widgets = widgets
}

// Focus set the focus on the Focusable Widget at the given index.
func (f *FocusManager) Focus(index int) {
	if index < 0 || index >= len(f.widgets) {
		return
	}

	if index != f.tabIndex {
		f.BlurCurrent()
	}

	f.tabIndex = index

	f.focus(f.tabIndex)
}

// FocusFirst gives the focus to the first active widget.
func (f *FocusManager) FocusFirst() {
	f.BlurCurrent()

	for i, item := range f.widgets {
		if item.IsActive() {
			f.tabIndex = i
			f.focus(f.tabIndex)
			return
		}
	}
}

// BlurCurrent simply blur the currently focused widget.
func (f *FocusManager) BlurCurrent() {
	if f.tabIndex >= 0 && f.tabIndex < len(f.widgets) {
		f.blur(f.tabIndex)
	}
}

// ForceInput forces the widget to enter input mode.
func (f *FocusManager) ForceInput(index int) {
	if f.tabIndex >= 0 && f.tabIndex < len(f.widgets) {
		f.enterInput(f.tabIndex)
	}
}

// ExitCurrentInput simply exits the currently inputting widget.
func (f *FocusManager) ExitCurrentInput() {
	if f.tabIndex >= 0 && f.tabIndex < len(f.widgets) {
		f.exitInput(f.tabIndex)
	}
}

// TabIndex returns the current tab index of the manager.
func (f *FocusManager) TabIndex() int {
	return f.tabIndex
}

// IsInputting returns true if a widget is in inputting mode.
func (f *FocusManager) IsInputting() bool {
	return f.isInputting
}

// PrevFocus moves the focus to the previous widget.
func (f *FocusManager) PrevFocus() {
	if f.widgets[f.tabIndex].IsFocused() {
		f.blur(f.tabIndex)
	}

	f.tabIndex = f.getPreviousIndex()

	f.focus(f.tabIndex)
}

// NextFocus moves the focus to the next widget.
func (f *FocusManager) NextFocus() {
	if f.widgets[f.tabIndex].IsFocused() {
		f.blur(f.tabIndex)
	}

	f.tabIndex = f.getNextIndex()

	f.focus(f.tabIndex)

}

// Update needs to be called in the screen or widget update function.
//
//	func (s *Screen) Update(msg tea.Msg) tea.Cmd {
//		switch msg := msg.(type) {
//		case tea.KeyMsg:
//			switch {
//			case key.Matches(msg, keybind.Esc):
//				return orvyn.SwitchToPreviousScreen()
//
//			case key.Matches(msg, keybind.Enter):
//				ok := s.submit()
//
//				if ok {
//					return orvyn.SwitchToPreviousScreen()
//				}
//
//				return nil
//			}
//		}
//
//		cmd := s.focusManager.Update(msg)
//
//		return cmd
//	}
func (f *FocusManager) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	cmd = nil

	if len(f.widgets) == 0 {
		return nil
	}

	if f.widgets[f.tabIndex].IsInputting() {
		var exitCmd tea.Cmd

		switch msg := msg.(type) {
		case tea.KeyMsg:
			if key.Matches(msg, f.widgets[f.tabIndex].GetExitInputKeybind()) {
				if f.widgets[f.tabIndex].CanExitInputting() {
					exitCmd = f.exitInput(f.tabIndex)
				}
			}
		}

		cmd = f.widgets[f.tabIndex].Update(msg)
		return tea.Batch(cmd, exitCmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, f.NextFocusKeybind):
			if f.ManageFocusNextPrevKeybind {
				f.NextFocus()
				return nil
			}

		case key.Matches(msg, f.PreviousFocusKeybind):
			if f.ManageFocusNextPrevKeybind {
				f.PrevFocus()
				return nil
			}
		}

		inputtingKeybind := f.widgets[f.tabIndex].GetEnterInputKeybind()

		if inputtingKeybind != nil {
			if key.Matches(msg, *inputtingKeybind) {
				cmd := f.enterInput(f.tabIndex)

				return cmd
			}
		}

		// Checking for specific focus keybind
		for i, widget := range f.widgets {
			keybind := widget.GetFocusKeybind()

			if keybind == nil {
				continue
			}

			if key.Matches(msg, *keybind) {
				if f.widgets[f.tabIndex].IsFocused() {
					f.blur(f.tabIndex)
				}

				f.tabIndex = i

				f.focus(f.tabIndex)

				return nil
			}
		}
	}

	// Call the update on the currently focused widget
	if f.widgets[f.tabIndex].IsFocused() {
		cmd = f.widgets[f.tabIndex].Update(msg)
	}

	return cmd
}

// Hidden functions

// focus is a shorthand to manage the focused state and call OnFocus.
func (f *FocusManager) focus(index int) {
	f.widgets[index].setFocused(true)
	f.widgets[index].OnFocus()
}

// blur is a shorthand to manage the focused state and call OnBlur.
func (f *FocusManager) blur(index int) {
	f.widgets[index].setFocused(false)
	f.widgets[index].OnBlur()
}

// enterInput is a shorthand to manage the inputting state and call OnEnterInput.
func (f *FocusManager) enterInput(index int) tea.Cmd {
	f.widgets[index].setInputting(true)
	cmd := f.widgets[index].OnEnterInput()
	f.isInputting = true

	return cmd
}

// exitInput is a shorthand to manage the inputting state and call OnExitInput.
func (f *FocusManager) exitInput(index int) tea.Cmd {
	f.widgets[index].setInputting(false)
	cmd := f.widgets[index].OnExitInput()
	f.isInputting = false

	return cmd
}

func (f *FocusManager) getNextIndex() int {
	var index int

	if len(f.widgets) == 0 {
		return 0
	}

	index = f.tabIndex + 1

	for {
		if index > len(f.widgets)-1 {
			index = 0
		}

		if !f.widgets[index].IsActive() {
			index++
			continue
		}

		return index
	}
}

func (f *FocusManager) getPreviousIndex() int {
	var index int

	if len(f.widgets) == 0 {
		return 0
	}

	index = f.tabIndex - 1

	for {
		if index < 0 {
			index = len(f.widgets) - 1
		}

		if !f.widgets[index].IsActive() {
			index--
			continue
		}

		return index
	}
}
