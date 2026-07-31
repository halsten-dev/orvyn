package layout

import (
	"math"
	"strings"

	"github.com/halsten-dev/orvyn"
)

// writeElements joins the rendered elements with a newline. An element that was
// squeezed to zero height and renders nothing is dropped entirely: keeping it
// would cost a blank line the layout has no room for. An element that renders
// nothing but still owns some height keeps its blank line, as it is spacing the
// layout asked for.
func writeElements(b *strings.Builder, elements []orvyn.Renderable) {
	first := true

	for _, e := range elements {
		view := e.Render()

		if view == "" && e.GetSize().Height <= 0 {
			continue
		}

		if !first {
			b.WriteString("\n")
		}

		b.WriteString(view)

		first = false
	}
}

// isFixedHeight reports whether an element keeps its own height instead of
// sharing the layout's leftover height. An element is fixed when it has no
// preferred height, or when its preferred height matches its minimal height.
func isFixedHeight(e orvyn.Renderable) bool {
	minHeight := e.GetMinSize().Height
	prefHeight := e.GetPreferredSize().Height

	return prefHeight == 0 || minHeight == prefHeight
}

// fixedHeight returns the height a fixed element should keep, falling back to
// its current size when it reports no minimal height.
func fixedHeight(e orvyn.Renderable) int {
	height := e.GetMinSize().Height

	if height == 0 {
		height = e.GetSize().Height
	}

	return height
}

// resizeFlexibleElements gives every fixed-height element its own height, then
// shares the remaining height among the flexible ones. Each flexible element is
// guaranteed its minimal height first, and only the surplus on top of it is split
// proportionally to the preferred heights. Allocating below the minimum would be
// a lie: the widget renders at its minimum anyway and the layout overflows by the
// difference.
//
// When there is not even room for every minimum, the minimums themselves are what
// gets shared, scaled down proportionally, so one element cannot starve while
// another keeps its full share.
//
// The rounding remainder goes to the last flexible element, so the allocation
// always adds up to availableHeight exactly.
//
// Elements must be the ones actually being rendered: an element that is not
// passed here reserves no height. The fixed/flexible split is recomputed on
// every call because widgets change their reported sizes at runtime.
func resizeFlexibleElements(width, availableHeight int, elements ...orvyn.Renderable) {
	flexible := make([]orvyn.Renderable, 0, len(elements))

	remaining := availableHeight

	for _, e := range elements {
		if isFixedHeight(e) {
			height := fixedHeight(e)

			e.Resize(orvyn.NewSize(width, height))
			remaining -= height

			continue
		}

		flexible = append(flexible, e)
	}

	if len(flexible) == 0 {
		return
	}

	remaining = max(remaining, 0)

	minTotal := 0
	prefTotal := 0

	for _, e := range flexible {
		minTotal += e.GetMinSize().Height
		prefTotal += e.GetPreferredSize().Height
	}

	// Every element starts at its minimum and the surplus is shared by preferred
	// height...
	guaranteed := func(e orvyn.Renderable) int { return e.GetMinSize().Height }
	weight := func(e orvyn.Renderable) int { return e.GetPreferredSize().Height }

	surplus := remaining - minTotal
	weightTotal := prefTotal

	// ...unless the minimums alone do not fit, in which case they are the thing
	// being shared.
	if surplus < 0 {
		guaranteed = func(orvyn.Renderable) int { return 0 }
		weight = func(e orvyn.Renderable) int { return e.GetMinSize().Height }

		surplus = remaining
		weightTotal = minTotal
	}

	left := remaining

	for i, e := range flexible {
		var height int

		switch {
		case i == len(flexible)-1:
			height = max(left, 0)
		case weightTotal == 0:
			height = guaranteed(e)
		default:
			height = guaranteed(e) + int(math.Round(
				float64(surplus)*float64(weight(e))/float64(weightTotal)))
		}

		e.Resize(orvyn.NewSize(width, height))
		left -= height
	}
}
