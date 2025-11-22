package layout

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/halsten-dev/orvyn"
)

// HBoxGrowLayout is a layout that uniformize the height of all the elements
// and calculate the width of each element.
type HBoxGrowLayout struct {
	orvyn.BaseLayout

	// gap defines the space between every element.
	gap int

	// compensatorIndex is the index of the widget that will be responsible to compensate
	// the possible calculation precision loss.
	compensatorIndex int

	// fullHeight takes the full layout height.
	fullHeight bool
}

// NewHBoxGrowLayout creates a new instance of a horizontal grow box layout.
// gap : defines the space between every element.
// compensatorIndex : specify which element of the layout will take the compensation size
// to match the layout size.
func NewHBoxGrowLayout(gap, compensatorIndex int, elements ...orvyn.Renderable) *HBoxGrowLayout {
	l := new(HBoxGrowLayout)

	l.BaseLayout = orvyn.NewBaseLayout(elements...)
	l.gap = gap
	l.compensatorIndex = compensatorIndex
	l.fullHeight = false

	return l
}

// NewHBoxGrowFullHeightLayout creates a new instance of a horizontal grow box layout that takes full height.
// gap : defines the space between every element.
// compensatorIndex : specify which element of the layout will take the compensation size
// to match the layout size.
func NewHBoxGrowFullHeightLayout(gap, compensatorIndex int, elements ...orvyn.Renderable) *HBoxGrowLayout {
	l := NewHBoxGrowLayout(gap, compensatorIndex, elements...)

	l.fullHeight = true

	return l
}

func (l *HBoxGrowLayout) Render() string {
	var view []string
	var elementSize orvyn.Size

	if len(l.GetElements()) == 0 {
		return ""
	}

	layoutSize := l.GetSize()

	elementSize.Height = layoutSize.Height

	if !l.fullHeight {
		minSize := l.GetMinSize()
		prefSize := l.GetPreferredSize()

		if layoutSize.Height <= minSize.Height {
			elementSize.Height = minSize.Height
		} else if layoutSize.Height >= prefSize.Height {
			elementSize.Height = prefSize.Height
		}
	}

	availableWidth := layoutSize.Width - (l.gap*len(l.GetElements()) - 1)
	elementSize.Width = int(math.Floor(float64(availableWidth / len(l.GetElements()))))

	// calculate the compensation
	compensatorSize := l.calculateCompensatorSize(elementSize, layoutSize)

	view = make([]string, 0)

	for i, e := range l.GetElements() {
		if i > 0 {
			view = append(view, strings.Repeat(" ", l.gap))
		}

		if i == l.compensatorIndex {
			e.Resize(compensatorSize)
		} else {
			e.Resize(elementSize)
		}

		view = append(view, e.Render())
	}

	return lipgloss.JoinHorizontal(lipgloss.Center,
		view...)
}

func (l *HBoxGrowLayout) calculateCompensatorSize(baseElementSize, layoutSize orvyn.Size) orvyn.Size {
	totalWidth := 0

	for i := range l.GetElements() {
		if i > 0 {
			totalWidth += l.gap
		}

		totalWidth += baseElementSize.Width
	}

	compensation := layoutSize.Width - totalWidth

	baseElementSize.Width += compensation

	return baseElementSize
}

func (l *HBoxGrowLayout) GetMinSize() orvyn.Size {
	var size orvyn.Size

	for _, e := range l.GetElements() {
		eSize := e.GetMinSize()

		if eSize.Width == 0 && eSize.Height == 0 {
			eSize = e.GetSize()
		}

		size.Width = max(eSize.Width, size.Width)
		size.Height = max(eSize.Height, size.Height)
	}

	size.Width *= len(l.GetElements())

	return size
}

func (l *HBoxGrowLayout) GetPreferredSize() orvyn.Size {
	var size orvyn.Size

	for _, e := range l.GetElements() {
		eSize := e.GetPreferredSize()

		if eSize.Width == 0 && eSize.Height == 0 {
			eSize = e.GetSize()
		}

		size.Width = max(eSize.Width, size.Width)
		size.Height = max(eSize.Height, size.Height)
	}

	size.Width *= len(l.GetElements())

	return size
}
