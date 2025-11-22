package orvyn

// Activable interface represents something that can be active or not.
type Activable interface {
	// SetActive defines the active status of the Activable.
	SetActive(bool)

	// IsActive returns the active status of the Activable.
	IsActive() bool
}

type BaseActivable struct {
	active bool
}

func (b *BaseActivable) SetActive(active bool) {
	b.active = active
}

func (b *BaseActivable) IsActive() bool {
	return b.active
}

func NewBaseActivable() BaseActivable {
	a := BaseActivable{}

	a.active = true

	return a
}
