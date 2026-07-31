package layout

import (
	"strings"

	"github.com/halsten-dev/orvyn"
)

// DefinedWidthVerticalLayout arranges elements vertically within the given width
// bounds. Elements whose minimal and preferred heights differ share the leftover
// height proportionally to their preferred heights; the others keep their own
// height. Widgets that return 0 or the same min and preferred height are
// considered fixed size widgets.
type DefinedWidthVerticalLayout struct {
	orvyn.BaseLayout

	PreferredWidth int
	MinWidth       int

	// Margin is the space kept free around the elements. Width is the total
	// horizontal space removed, Height the total vertical one.
	Margin orvyn.Size
}

func NewDefinedWidthVerticalLayout(minWidth int, prefWidth int, margin orvyn.Size, elements ...orvyn.Renderable) *DefinedWidthVerticalLayout {
	l := new(DefinedWidthVerticalLayout)

	l.BaseLayout = orvyn.NewBaseLayout(elements...)

	l.MinWidth = minWidth
	l.PreferredWidth = prefWidth
	l.Margin = margin

	return l
}

func (l *DefinedWidthVerticalLayout) Render() string {
	var b strings.Builder

	visibleElements := l.GetElements()

	if len(visibleElements) == 0 {
		return ""
	}

	layoutSize := l.GetSize()

	resizeFlexibleElements(
		l.fitWidth(layoutSize.Width),
		max(layoutSize.Height-l.Margin.Height, 0),
		visibleElements...)

	writeElements(&b, visibleElements)

	return b.String()
}

// fitWidth returns the element width: the available width minus the margin, kept
// within the layout min and preferred widths. The available width always wins so
// a terminal narrower than MinWidth shrinks the content instead of overflowing.
func (l *DefinedWidthVerticalLayout) fitWidth(layoutWidth int) int {
	width := layoutWidth - l.Margin.Width

	width = min(width, l.PreferredWidth-l.Margin.Width)
	width = max(width, l.MinWidth-l.Margin.Width)
	width = min(width, layoutWidth)

	return max(width, 0)
}

func (l *DefinedWidthVerticalLayout) GetMinSize() orvyn.Size {
	size := orvyn.NewSize(l.MinWidth, l.Margin.Height)

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
	size := orvyn.NewSize(l.PreferredWidth, l.Margin.Height)

	for _, e := range l.GetElements() {
		height := e.GetPreferredSize().Height

		if height == 0 {
			height = e.GetSize().Height
		}

		size.Height += height
	}

	return size
}
