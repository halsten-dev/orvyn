package orvyn

// Activable interface represents something that can be active or not.
type Activable interface {
	// SetActive defines the active status of the Activable.
	SetActive(bool)

	// IsActive returns the active status of the Activable.
	IsActive() bool
}

// BaseActivable interface represents the basic implementation of an Activable.
// To avoid code repetition.
type BaseActivable struct {
	active bool
}

func (b *BaseActivable) SetActive(active bool) {
	b.active = active
}

func (b *BaseActivable) IsActive() bool {
	return b.active
}

// NewBaseActivable creates a new BaseActivable.
func NewBaseActivable() BaseActivable {
	a := BaseActivable{}

	a.active = true

	return a
}
