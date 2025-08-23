package orvyn

type Activable interface {
	SetActive(bool)
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
