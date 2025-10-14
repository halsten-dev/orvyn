package orvyn

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/halsten-dev/orvyn/theme"
)

type Widget interface {
	// Init is called on the Widget when entering a Screen that contains it.
	Init() tea.Cmd

	// SetStyle allows to set the style of the base widget. By default theme.BlurredWidgetStyleID.
	SetStyle(lipgloss.Style)

	// GetStyle returns the widget style for custom rendering.
	GetStyle() lipgloss.Style

	// GetContentSize must be used when rendering the widget to get the real available size for the widgets content.
	// Borders of the style have been taken into account.
	GetContentSize() Size

	// Updatable Widget can be updated.
	Updatable

	// Renderable Widget can be rendered.
	Renderable
}

type BaseWidget struct {
	BaseRenderable
	style       lipgloss.Style
	contentSize Size
}

func NewBaseWidget() BaseWidget {
	w := BaseWidget{}

	w.BaseRenderable = NewBaseRenderable()
	w.style = GetTheme().Style(theme.BlurredWidgetStyleID)

	return w
}

func (b *BaseWidget) Init() tea.Cmd {
	return nil
}

func (b *BaseWidget) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (b *BaseWidget) Resize(size Size) {
	b.BaseRenderable.Resize(size)

	size.Width -= b.style.GetHorizontalFrameSize()
	size.Height -= b.style.GetVerticalFrameSize()

	size.Width = max(size.Width, 0)
	size.Height = max(size.Height, 0)

	b.contentSize = size
}

func (b *BaseWidget) GetContentSize() Size {
	return b.contentSize
}

func (b *BaseWidget) SetStyle(style lipgloss.Style) {
	b.style = style
}

func (b *BaseWidget) GetStyle() lipgloss.Style {
	return b.style
}
