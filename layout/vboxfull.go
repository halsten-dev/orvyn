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

	autoFlex bool
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
		resizeFlexibleElements(width, max(layoutSize.Height-l.margin.Height, 0), visibleElements...)
	} else {
		l.resizeSingleGrow(width, layoutSize)
	}

	writeElements(&b, visibleElements)

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
