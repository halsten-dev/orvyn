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
// shares the remaining height among the flexible ones proportionally to their
// preferred heights. The rounding remainder goes to the last flexible element so
// the allocation always adds up to availableHeight exactly.
//
// Elements must be the ones actually being rendered: an element that is not
// passed here reserves no height. The fixed/flexible split is recomputed on
// every call because widgets change their reported sizes at runtime.
func resizeFlexibleElements(width, availableHeight int, elements ...orvyn.Renderable) {
	flexible := make([]orvyn.Renderable, 0, len(elements))

	remaining := availableHeight
	totalFlexPref := 0

	for _, e := range elements {
		if isFixedHeight(e) {
			height := fixedHeight(e)

			e.Resize(orvyn.NewSize(width, height))
			remaining -= height

			continue
		}

		flexible = append(flexible, e)
		totalFlexPref += e.GetPreferredSize().Height
	}

	remaining = max(remaining, 0)

	if len(flexible) == 0 || totalFlexPref == 0 {
		return
	}

	left := remaining

	for i, e := range flexible {
		var height int

		if i == len(flexible)-1 {
			height = max(left, 0)
		} else {
			height = int(math.Round(
				float64(remaining) * float64(e.GetPreferredSize().Height) / float64(totalFlexPref)))
		}

		e.Resize(orvyn.NewSize(width, height))
		left -= height
	}
}
