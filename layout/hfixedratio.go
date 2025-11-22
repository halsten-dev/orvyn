package layout

import (
	"log"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/halsten-dev/orvyn"
)

type FixedRatioRenderable struct {
	ratio     float64
	element   orvyn.Renderable
	tempWidth int
}

func NewFixedRatioRenderable(ratio float64, element orvyn.Renderable) FixedRatioRenderable {
	return FixedRatioRenderable{
		ratio:   ratio,
		element: element,
	}
}

type HBoxFixedRatio struct {
	orvyn.BaseLayout

	elements []FixedRatioRenderable

	margin int

	gap int

	compensatorIndex int
}

func NewHBoxFixedRatioLayout(margin int, gap,
	compensatorIndex int, elements ...FixedRatioRenderable) *HBoxFixedRatio {
	l := new(HBoxFixedRatio)

	l.BaseLayout = orvyn.NewBaseLayout()

	l.margin = margin
	l.gap = gap
	l.compensatorIndex = compensatorIndex
	l.elements = elements

	totalRatio := 0.0

	for _, e := range l.elements {
		totalRatio += e.ratio
	}

	if totalRatio != 1 {
		log.Fatal("HBoxFixedRatioLayout : Total elements ratio not equals 1.")
	}

	return l
}

func (l *HBoxFixedRatio) Render() string {
	layoutSize := l.GetSize()

	totalWidth := layoutSize.Width - l.margin

	if l.gap > 0 {
		totalWidth -= l.gap*len(l.elements) - 1
	}

	usedWidth := 0

	elementSize := orvyn.NewSize(0, layoutSize.Height)

	for i, e := range l.elements {
		width := int(math.Floor(float64(totalWidth) * e.ratio))
		l.elements[i].tempWidth = width
		usedWidth += width
	}

	compensationWidth := totalWidth - usedWidth

	l.elements[l.compensatorIndex].tempWidth += compensationWidth

	view := make([]string, 0)

	for i, e := range l.elements {
		if i > 0 {
			view = append(view, strings.Repeat(" ", l.gap))
		}

		elementSize.Width = e.tempWidth

		e.element.Resize(elementSize)

		view = append(view, e.element.Render())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		view...)
}

func (l *HBoxFixedRatio) GetMinSize() orvyn.Size {
	var size orvyn.Size

	for _, e := range l.elements {
		eSize := e.element.GetMinSize()

		if eSize.Width == 0 && eSize.Height == 0 {
			eSize = e.element.GetSize()
		}

		size.Width = max(eSize.Width, size.Width)
		size.Height = max(eSize.Height, size.Height)
	}

	size.Width *= len(l.elements)

	return size
}

func (l *HBoxFixedRatio) GetPreferredSize() orvyn.Size {
	var size orvyn.Size

	for _, e := range l.elements {
		eSize := e.element.GetPreferredSize()

		if eSize.Width == 0 && eSize.Height == 0 {
			eSize = e.element.GetSize()
		}

		size.Width = max(eSize.Width, size.Width)
		size.Height = max(eSize.Height, size.Height)
	}

	size.Width *= len(l.elements)

	return size
}
