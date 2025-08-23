package orvyn

type Layout interface {
	Renderable

	GetElements() []Renderable
}

type BaseLayout struct {
	BaseRenderable

	elements []Renderable
}

func NewBaseLayout(elements []Renderable) BaseLayout {
	b := BaseLayout{}

	b.BaseRenderable = NewBaseRenderable()
	b.elements = elements

	return b
}

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

func (b *BaseLayout) SetActive(active bool) {
	for _, e := range b.elements {
		e.SetActive(active)
	}

	b.active = active
}

func (b *BaseLayout) IsActive() bool {
	return b.active
}
