package layout

import (
	"math"
	"strings"

	"github.com/halsten-dev/orvyn"
)

type VBoxFullLayout struct {
	orvyn.BaseLayout

	margin     orvyn.Size
	growWidget orvyn.Renderable
	maxWidth   bool

	autoFlex               bool
	fixedHeightElements    []orvyn.Renderable
	flexibleHeightElements []flexibleHeightElement
}

func NewVBoxFullLayout(margin orvyn.Size, growIndex int, elements ...orvyn.Renderable) *VBoxFullLayout {
	l := new(VBoxFullLayout)

	l.BaseLayout = orvyn.NewBaseLayout(elements...)
	l.growWidget = elements[growIndex]
	l.maxWidth = false
	l.margin = margin

	return l
}

func NewMaxWidthVBoxFullLayout(margin orvyn.Size, growIndex int, elements ...orvyn.Renderable) *VBoxFullLayout {
	l := new(VBoxFullLayout)

	l.BaseLayout = orvyn.NewBaseLayout(elements...)
	l.growWidget = elements[growIndex]
	l.maxWidth = true
	l.margin = margin

	return l
}

// NewFlexibleVBoxFullLayout arranges elements vertically and shares the leftover
// height among every flexible-height element, proportionally to their preferred
// heights. An element is flexible when GetMinSize().Height < GetPreferredSize().Height;
// elements whose preferred height is 0 or equal to their min height are fixed.
func NewFlexibleVBoxFullLayout(margin orvyn.Size, elements ...orvyn.Renderable) *VBoxFullLayout {
	return newFlexibleVBoxFullLayout(margin, false, elements...)
}

// NewMaxWidthFlexibleVBoxFullLayout behaves like NewFlexibleVBoxFullLayout but
// ignores the width min/preferred constraints (full width).
func NewMaxWidthFlexibleVBoxFullLayout(margin orvyn.Size, elements ...orvyn.Renderable) *VBoxFullLayout {
	return newFlexibleVBoxFullLayout(margin, true, elements...)
}

func newFlexibleVBoxFullLayout(margin orvyn.Size, maxWidth bool, elements ...orvyn.Renderable) *VBoxFullLayout {
	l := new(VBoxFullLayout)

	l.BaseLayout = orvyn.NewBaseLayout(elements...)
	l.margin = margin
	l.maxWidth = maxWidth
	l.autoFlex = true

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

		l.flexibleHeightElements = append(l.flexibleHeightElements, flexibleHeightElement{
			element: e,
		})
	}

	return l
}

func (l *VBoxFullLayout) Render() string {
	var b strings.Builder

	visibleElements := l.GetElements()

	if len(visibleElements) == 0 {
		return ""
	}

	layoutSize := l.GetSize()
	width := l.fitWidth(layoutSize.Width)

	if l.autoFlex {
		l.resizeAutoFlex(width, layoutSize.Height)
	} else {
		l.resizeSingleGrow(width, layoutSize)
	}

	for i, e := range visibleElements {
		if i > 0 {
			b.WriteString("\n")
		}

		b.WriteString(e.Render())
	}

	return b.String()
}

// fitWidth returns the element width after applying margin and, unless maxWidth is
// set, the layout min/preferred width constraints.
func (l *VBoxFullLayout) fitWidth(layoutWidth int) int {
	width := layoutWidth - l.margin.Width

	if l.maxWidth {
		return width
	}

	minSize := l.GetMinSize()
	prefSize := l.GetPreferredSize()

	if width <= minSize.Width {
		width = minSize.Width - l.margin.Width
	} else if width >= prefSize.Width {
		width = prefSize.Width - l.margin.Width
	}

	return width
}

// resizeSingleGrow resizes elements for the classic single grow-widget behavior.
func (l *VBoxFullLayout) resizeSingleGrow(width int, layoutSize orvyn.Size) {
	for _, e := range l.GetElements() {
		if e == l.growWidget {
			continue
		}

		e.Resize(orvyn.NewSize(width, e.GetMinSize().Height))
	}

	if l.growWidget != nil {
		l.growWidget.Resize(l.calculateGrowSize(orvyn.NewSize(width, 0), layoutSize))
	}
}

// resizeAutoFlex resizes fixed elements to their own height then shares the
// leftover height among flexible elements proportionally to their preferred height.
// Rounding remainder goes to the last flexible element so the layout fills its height.
func (l *VBoxFullLayout) resizeAutoFlex(width, layoutHeight int) {
	remaining := layoutHeight - l.margin.Height

	for _, e := range l.fixedHeightElements {
		height := e.GetMinSize().Height

		if height == 0 {
			height = e.GetSize().Height
		}

		e.Resize(orvyn.NewSize(width, height))
		remaining -= height
	}

	remaining = max(remaining, 0)

	totalFlexPref := 0
	for _, e := range l.flexibleHeightElements {
		totalFlexPref += e.element.GetPreferredSize().Height
	}

	if totalFlexPref == 0 {
		return
	}

	left := remaining

	for i, e := range l.flexibleHeightElements {
		var height int

		if i == len(l.flexibleHeightElements)-1 {
			height = max(left, 0)
		} else {
			height = int(math.Round(
				float64(remaining) * float64(e.element.GetPreferredSize().Height) / float64(totalFlexPref)))
		}

		e.element.Resize(orvyn.NewSize(width, height))
		left -= height
	}
}

func (l *VBoxFullLayout) calculateGrowSize(elementSize, layoutSize orvyn.Size) orvyn.Size {
	totalHeight := layoutSize.Height

	for _, e := range l.GetElements() {
		if e == l.growWidget {
			continue
		}

		height := e.GetSize().Height

		if height == 0 {
			height = e.GetMinSize().Height
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
