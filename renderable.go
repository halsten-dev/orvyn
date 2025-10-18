package orvyn

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Renderable interface {
	Activable

	Render() string
	Resize(Size)
	GetSize() Size
	SetMinSize(Size)
	GetMinSize() Size
	SetPreferredSize(Size)
	GetPreferredSize() Size
}

type BaseRenderable struct {
	BaseActivable

	size          Size
	minSize       Size
	preferredSize Size
}

func NewBaseRenderable() BaseRenderable {
	b := BaseRenderable{}

	b.BaseActivable = NewBaseActivable()

	b.minSize = NewSize(1, 1)
	b.preferredSize = NewSize(1, 1)

	return b
}

func (b *BaseRenderable) Resize(size Size) {
	b.size = size
}

func (b *BaseRenderable) GetSize() Size {
	return b.size
}

func (b *BaseRenderable) SetMinSize(size Size) {
	b.minSize = size
}

func (b *BaseRenderable) GetMinSize() Size {
	return b.minSize
}

func (b *BaseRenderable) SetPreferredSize(size Size) {
	b.preferredSize = size
}

func (b *BaseRenderable) GetPreferredSize() Size {
	return b.preferredSize
}

type Updatable interface {
	Update(tea.Msg) tea.Cmd
}
