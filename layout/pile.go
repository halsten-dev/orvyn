package layout

import (
	"github.com/halsten-dev/orvyn"
)

// PileLayout Show only the first visible widget in the full size available.
type PileLayout struct {
	orvyn.BaseLayout
}

func NewPileLayout(elements ...orvyn.Renderable) *PileLayout {
	l := new(PileLayout)

	l.BaseLayout = orvyn.NewBaseLayout(elements...)

	return l
}

func (l *PileLayout) Render() string {
	var view string

	if len(l.GetElements()) == 0 {
		return ""
	}

	layoutSize := l.GetSize()

	for _, e := range l.GetElements() {
		e.Resize(layoutSize)
		view = e.Render()
	}

	return view
}

func (l *PileLayout) GetMinSize() orvyn.Size {
	for _, e := range l.GetElements() {
		return e.GetMinSize()
	}

	return orvyn.NewSize(0, 0)
}

func (l *PileLayout) GetPreferredSize() orvyn.Size {
	for _, e := range l.GetElements() {
		return e.GetPreferredSize()
	}

	return orvyn.NewSize(0, 0)

}
