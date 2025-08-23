package orvyn

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Renderable interface {
	Activable

	Render() string
	Resize(Size)
	GetSize() Size
	GetMinSize() Size
	GetPreferredSize() Size
}

type BaseRenderable struct {
	BaseActivable

	size Size
}

func NewBaseRenderable() BaseRenderable {
	b := BaseRenderable{}

	b.BaseActivable = NewBaseActivable()

	return b
}

func (b *BaseRenderable) Resize(size Size) {
	b.size = size
}

func (b *BaseRenderable) GetSize() Size {
	return b.size
}

func (b *BaseRenderable) GetMinSize() Size {
	return NewSize(1, 1)
}

func (b *BaseRenderable) GetPreferredSize() Size {
	return NewSize(1, 1)
}

type Updatable interface {
	Update(tea.Msg) tea.Cmd
}
