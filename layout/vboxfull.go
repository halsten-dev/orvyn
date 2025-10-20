package layout

import (
	"strings"

	"github.com/halsten-dev/orvyn"
)

type VBoxFullLayout struct {
	orvyn.BaseLayout

	margin     orvyn.Size
	growWidget orvyn.Renderable
	maxWidth   bool
}

func NewVBoxFullLayout(margin orvyn.Size, growIndex int, elements []orvyn.Renderable) *VBoxFullLayout {
	l := new(VBoxFullLayout)

	l.BaseLayout = orvyn.NewBaseLayout(elements)
	l.growWidget = elements[growIndex]
	l.maxWidth = false
	l.margin = margin

	return l
}

func NewMaxWidthVBoxFullLayout(margin orvyn.Size, growIndex int, elements []orvyn.Renderable) *VBoxFullLayout {
	l := new(VBoxFullLayout)

	l.BaseLayout = orvyn.NewBaseLayout(elements)
	l.growWidget = elements[growIndex]
	l.maxWidth = true
	l.margin = margin

	return l
}

func (l *VBoxFullLayout) Render() string {
	var b strings.Builder
	var elementSize orvyn.Size
	var minSize orvyn.Size
	var prefSize orvyn.Size

	visibleElements := l.GetElements()

	if len(visibleElements) == 0 {
		return ""
	}

	layoutSize := l.GetSize()

	elementSize = orvyn.NewSize(layoutSize.Width-l.margin.Width, 0)

	if !l.maxWidth {
		minSize = l.GetMinSize()
		prefSize = l.GetPreferredSize()

		if elementSize.Width <= minSize.Width {
			elementSize.Width = minSize.Width - l.margin.Width
		} else if elementSize.Width >= prefSize.Width {
			elementSize.Width = prefSize.Width - l.margin.Width
		}
	}

	for _, e := range visibleElements {
		if e != l.growWidget {
			continue
		}

		elementSize.Height = e.GetMinSize().Height

		e.Resize(elementSize)
	}

	for i, e := range visibleElements {
		if i > 0 {
			b.WriteString("\n")
		}

		if e == l.growWidget {
			e.Resize(l.calculateGrowSize(elementSize, layoutSize))
		}

		b.WriteString(e.Render())
	}

	return b.String()
}

func (l *VBoxFullLayout) calculateGrowSize(elementSize, layoutSize orvyn.Size) orvyn.Size {
	totalHeight := layoutSize.Height

	for _, e := range l.GetElements() {
		if e == l.growWidget {
			continue
		}

		height := e.GetMinSize().Height

		if height == 0 {
			height = e.GetSize().Height
		}

		totalHeight -= height
	}

	totalHeight -= l.margin.Height

	return orvyn.NewSize(elementSize.Width, totalHeight)
}

func (l *VBoxFullLayout) GetMinSize() orvyn.Size {
	var size orvyn.Size

	for _, e := range l.GetElements() {
		eSize := e.GetMinSize()
		size.Height += eSize.Height

		size.Width = max(size.Width, eSize.Width)
	}

	return size
}

func (l *VBoxFullLayout) GetPreferredSize() orvyn.Size {
	var size orvyn.Size

	for _, e := range l.GetElements() {
		eSize := e.GetPreferredSize()
		size.Height += eSize.Height

		size.Width = max(size.Width, eSize.Width)
	}

	return size
}
