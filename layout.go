package orvyn

// Layout interface based on Renderable interface.
type Layout interface {
	// Renderable composing the Layout interface.
	Renderable

	// GetElements returns a slice of every active Renderable of the layout.
	GetElements() []Renderable
}

// BaseLayout type is used to simplify the creation of custom layouts.
type BaseLayout struct {
	BaseRenderable

	elements []Renderable
}

// NewBaseLayout creates and returns a new BaseLayout.
func NewBaseLayout(elements ...Renderable) BaseLayout {
	b := BaseLayout{}

	b.BaseRenderable = NewBaseRenderable()
	b.elements = elements

	return b
}

// GetElements returns a slice of activated []Renderable.
func (b *BaseLayout) GetElements() []Renderable {
	var visibleElements []Renderable

	for _, e := range b.elements {
		if !e.IsActive() {
			continue
		}

		visibleElements = append(visibleElements, e)
	}

	return visibleElements
}

// SetActive change the active state of all elements of the layout and the layout itself.
func (b *BaseLayout) SetActive(active bool) {
	for _, e := range b.elements {
		e.SetActive(active)
	}

	b.active = active
}
