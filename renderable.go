package orvyn

// Renderable interface defines any orvyn element that can be rendered.
type Renderable interface {
	// Activable compose the Renderable with the Activable interface
	Activable

	// Render is the way the renderable should be rendered.
	Render() string

	// Resize will be called when the renderable is resized by a layout.
	Resize(Size)

	// GetSize returns the current size of the Renderable.
	GetSize() Size

	// SetMinSize allows to set the minimal size of the Renderable.
	SetMinSize(Size)

	// GetMinSize returns the minimal size of the Renderable.
	GetMinSize() Size

	// SetPreferredSize allows to set the preferred size of the Renderable.
	SetPreferredSize(Size)

	// GetPreferredSize returns the preferred size of the Renderable.
	GetPreferredSize() Size
}

// BaseRenderable is usefull to bring default implementation of the Renderable interface.
type BaseRenderable struct {
	BaseActivable

	size          Size
	minSize       Size
	preferredSize Size
}

// NewBaseRenderable creates and returns a new BaseRenderable.
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
