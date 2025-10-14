package orvyn

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"github.com/halsten-dev/orvyn/theme"
)

type Focusable interface {
	// Updatable to be able to update a focusable widget.
	Updatable

	// Activable to be able to activate or deactivate widget.
	Activable

	// OnFocus is called when the widget gains the focus.
	OnFocus()

	// OnBlur is called when the widget is loosing focus.
	OnBlur()

	// OnEnterInput is called when the widget enters the input mode.
	// Input mode means that all the tea.Msg will be managed by the widget.
	OnEnterInput()

	// OnExitInput is called when the widget exits the input mode.
	OnExitInput()

	// IsFocused return true if the widget is currently focused.
	// A widget can be focused without being in input mode.
	IsFocused() bool

	// IsInputting return true if the widget is in input mode.
	// Input mode means that all the tea.Msg will be managed by the widget.
	IsInputting() bool

	// GetFocusKeybind can return a *key.Binding if the widget should be
	// able to get directly focused with one key.
	GetFocusKeybind() *key.Binding

	// GetEnterInputKeybind can return a *key.Binding if the widget should be
	// able to enter the input mode.
	// Input mode means that all the tea.Msg will be managed by the widget.
	GetEnterInputKeybind() *key.Binding

	// GetExitInputKeybind returns by default "Esc". Can be override.
	GetExitInputKeybind() key.Binding

	// CanExitInputting returns true if the widget can exit his inputting state.
	CanExitInputting() bool

	// setFocused allows to set the focused value.
	setFocused(bool)

	// setInputting allows to set the inputting value.
	setInputting(bool)

	// SetFocusedStyle allows to set the style when the widget is focused. By default theme.FocusedWidgetStyleID.
	SetFocusedStyle(lipgloss.Style)

	// SetBlurredStyle allows to set the style when the widget is blurred. By default theme.BlurredWidgetStyleID.
	SetBlurredStyle(lipgloss.Style)
}

// BaseFocusable can be integrated to a Widget to make it focusable.
type BaseFocusable struct {
	widget    Widget
	focused   bool
	inputting bool

	focusedStyle lipgloss.Style
	blurredStyle lipgloss.Style
}

func NewBaseFocusable(widget Widget) BaseFocusable {
	t := GetTheme()

	return BaseFocusable{
		widget:       widget,
		focused:      false,
		inputting:    false,
		focusedStyle: t.Style(theme.FocusedWidgetStyleID),
		blurredStyle: t.Style(theme.BlurredWidgetStyleID),
	}
}

func (b *BaseFocusable) OnFocus() {
	b.widget.SetStyle(b.focusedStyle)
}

func (b *BaseFocusable) OnBlur() {
	b.widget.SetStyle(b.blurredStyle)
}

func (b *BaseFocusable) IsFocused() bool {
	return b.focused
}

func (b *BaseFocusable) IsInputting() bool {
	return b.inputting
}

func (b *BaseFocusable) GetFocusKeybind() *key.Binding {
	return nil
}

func (b *BaseFocusable) GetEnterInputKeybind() *key.Binding {
	return nil
}

func (b *BaseFocusable) GetExitInputKeybind() key.Binding {
	return key.NewBinding(key.WithKeys("esc"))
}

func (b *BaseFocusable) CanExitInputting() bool {
	return true
}

func (b *BaseFocusable) setFocused(focused bool) {
	b.focused = focused
}

func (b *BaseFocusable) setInputting(input bool) {
	b.inputting = input
}

func (b *BaseFocusable) SetFocusedStyle(style lipgloss.Style) {
	b.focusedStyle = style
}

func (b *BaseFocusable) SetBlurredStyle(style lipgloss.Style) {
	b.blurredStyle = style
}
