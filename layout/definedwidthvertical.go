package layout

import (
	"math"
	"strings"

	"github.com/halsten-dev/orvyn"
)

type flexibleHeightElement struct {
	element     orvyn.Renderable
	heightRatio float64
}

// DefinedWidthVerticalLayout arranges vertically elements within the given width values.
// The height of each element will be defined between the min and the preferred size.
// Widgets that return 0,0 or the same min and preferred size are considered as fixed size widgets.
type DefinedWidthVerticalLayout struct {
	orvyn.BaseLayout

	PreferredWidth int
	MinWidth       int
	Margin         int

	fixedHeightElements    []orvyn.Renderable
	flexibleHeightElements []flexibleHeightElement
}

func NewDefinedWidthVerticalLayout(minWidth int, prefWidth int, margin int, elements ...orvyn.Renderable) *DefinedWidthVerticalLayout {
	l := new(DefinedWidthVerticalLayout)

	l.BaseLayout = orvyn.NewBaseLayout(elements...)

	l.MinWidth = minWidth
	l.PreferredWidth = prefWidth
	l.Margin = margin

	l.fixedHeightElements = make([]orvyn.Renderable, 0)
	l.flexibleHeightElements = make([]flexibleHeightElement, 0)

	for _, e := range elements {
		fixedHeight := false

		eMinHeight := e.GetMinSize().Height
		ePrefHeight := e.GetPreferredSize().Height

		switch {
		case ePrefHeight == 0:
			fixedHeight = true
		case eMinHeight == ePrefHeight:
			fixedHeight = true
		}

		if fixedHeight {
			l.fixedHeightElements = append(l.fixedHeightElements, e)
			continue
		}

		el := flexibleHeightElement{
			element: e,
		}

		l.flexibleHeightElements = append(l.flexibleHeightElements, el)
	}

	return l
}

func (l *DefinedWidthVerticalLayout) Render() string {
	var b strings.Builder
	var s orvyn.Size
	var minSize orvyn.Size
	var prefSize orvyn.Size

	if len(l.GetElements()) == 0 {
		return ""
	}

	l.calculateHeightRatio()

	size := l.GetSize()

	s = orvyn.NewSize(size.Width-l.Margin, size.Height-l.Margin)

	minSize = l.GetMinSize()
	prefSize = l.GetPreferredSize()

	if s.Width <= minSize.Width {
		s.Width = minSize.Width - l.Margin
	} else if s.Width >= prefSize.Width {
		s.Width = prefSize.Width - l.Margin
	}

	s.Height = max(s.Height, 0)

	l.calculateSize(s)

	for i, e := range l.GetElements() {
		if i > 0 {
			b.WriteString("\n")
		}

		b.WriteString(e.Render())
	}

	return b.String()
}

func (l *DefinedWidthVerticalLayout) calculateSize(size orvyn.Size) {
	totalHeight := size.Height

	for _, e := range l.fixedHeightElements {
		eHeight := e.GetSize().Height
		totalHeight -= eHeight
		e.Resize(orvyn.NewSize(size.Width, eHeight))
	}

	for _, e := range l.flexibleHeightElements {
		eHeight := int(math.Round(float64(totalHeight) * e.heightRatio))

		e.element.Resize(orvyn.NewSize(size.Width, eHeight))
	}
}

func (l *DefinedWidthVerticalLayout) GetMinSize() orvyn.Size {
	var size orvyn.Size

	size.Width = l.MinWidth

	for _, e := range l.GetElements() {
		height := e.GetMinSize().Height

		if height == 0 {
			height = e.GetSize().Height
		}

		size.Height += height
	}

	return size
}

func (l *DefinedWidthVerticalLayout) GetPreferredSize() orvyn.Size {
	var size orvyn.Size

	size.Width = l.PreferredWidth

	for _, e := range l.GetElements() {
		height := e.GetPreferredSize().Height

		if height == 0 {
			height = e.GetSize().Height
		}

		size.Height += height
	}

	return size
}

func (l *DefinedWidthVerticalLayout) calculateHeightRatio() {
	maxHeight := l.GetPreferredSize().Height

	fixedElementHeight := 0

	for _, e := range l.fixedHeightElements {
		height := e.GetPreferredSize().Height

		if height == 0 {
			height = e.GetSize().Height
		}

		fixedElementHeight -= height
	}

	maxHeight -= fixedElementHeight

	for i, e := range l.flexibleHeightElements {
		l.flexibleHeightElements[i].heightRatio =
			float64(e.element.GetPreferredSize().Height) / float64(maxHeight)
	}
}
