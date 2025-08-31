package layout

import (
	"github.com/halsten-dev/orvyn"
	"strings"
)

// VBoxLayout arranges elements vertically with flexible width (base on the largest width).
// All heights are the minimal for each element.
type VBoxLayout struct {
	orvyn.BaseLayout

	margin int

	// maxWidth defines if element should take the whole available width.
	maxWidth bool
}

func NewVBoxLayout(margin int, elements []orvyn.Renderable) *VBoxLayout {
	l := new(VBoxLayout)

	l.BaseLayout = orvyn.NewBaseLayout(elements)
	l.margin = margin
	l.maxWidth = false

	return l
}

func NewMaxWidthVBoxLayout(margin int, elements []orvyn.Renderable) *VBoxLayout {
	l := new(VBoxLayout)

	l.BaseLayout = orvyn.NewBaseLayout(elements)
	l.margin = margin
	l.maxWidth = true

	return l
}

func (l *VBoxLayout) Render() string {
	var b strings.Builder
	var s orvyn.Size
	var minSize orvyn.Size
	var prefSize orvyn.Size

	if len(l.GetElements()) == 0 {
		return ""
	}

	layoutSize := l.GetSize()

	s = orvyn.NewSize(layoutSize.Width-l.margin, 0)

	if !l.maxWidth {
		minSize = l.GetMinSize()
		prefSize = l.GetPreferredSize()

		if s.Width <= minSize.Width {
			s.Width = minSize.Width - l.margin
		} else if s.Width >= prefSize.Width {
			s.Width = prefSize.Width - l.margin
		}
	}

	for i, e := range l.GetElements() {
		if i > 0 {
			b.WriteString("\n")
		}

		s.Height = e.GetMinSize().Height

		e.Resize(s)
		b.WriteString(e.Render())
	}

	return b.String()
}

func (l *VBoxLayout) GetMinSize() orvyn.Size {
	var size orvyn.Size

	for _, e := range l.GetElements() {
		eSize := e.GetMinSize()
		size.Height += eSize.Height

		size.Width = max(size.Width, eSize.Width)
	}

	return size
}

func (l *VBoxLayout) GetPreferredSize() orvyn.Size {
	var size orvyn.Size

	for _, e := range l.GetElements() {
		eSize := e.GetPreferredSize()
		size.Height += eSize.Height

		size.Width = max(size.Width, eSize.Width)
	}

	return size
}
