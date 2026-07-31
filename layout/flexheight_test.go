package layout

import (
	"strings"
	"testing"

	"github.com/halsten-dev/orvyn"
)

// flexStub is a minimal Renderable whose min/preferred sizes can be changed at
// runtime, so tests can check how layouts react to widgets that resize themselves.
type flexStub struct {
	orvyn.BaseRenderable

	min  orvyn.Size
	pref orvyn.Size
}

func newFlexStub(min, pref orvyn.Size) *flexStub {
	s := new(flexStub)

	s.BaseRenderable = orvyn.NewBaseRenderable()
	s.min = min
	s.pref = pref

	return s
}

func (s *flexStub) GetMinSize() orvyn.Size       { return s.min }
func (s *flexStub) GetPreferredSize() orvyn.Size { return s.pref }

// Render draws exactly the height it was allocated so tests can measure overflow.
func (s *flexStub) Render() string {
	size := s.GetSize()

	if size.Height <= 0 {
		return ""
	}

	lines := make([]string, 0, size.Height)

	for range size.Height {
		lines = append(lines, strings.Repeat("#", max(size.Width, 0)))
	}

	return strings.Join(lines, "\n")
}

func allocatedHeights(elements ...orvyn.Renderable) []int {
	heights := make([]int, 0, len(elements))

	for _, e := range elements {
		heights = append(heights, e.GetSize().Height)
	}

	return heights
}

func sum(values []int) int {
	total := 0

	for _, v := range values {
		total += v
	}

	return total
}

// A fixed element (min == preferred) must keep its own height, and that height
// must be removed from the pool the flexible elements share.
func TestResizeFlexibleElementsSubtractsFixedHeights(t *testing.T) {
	fixed := newFlexStub(orvyn.NewSize(10, 3), orvyn.NewSize(10, 3))
	flexA := newFlexStub(orvyn.NewSize(10, 5), orvyn.NewSize(10, 10))
	flexB := newFlexStub(orvyn.NewSize(10, 5), orvyn.NewSize(10, 30))

	resizeFlexibleElements(40, 23, fixed, flexA, flexB)

	if got := fixed.GetSize().Height; got != 3 {
		t.Errorf("fixed element height = %d, want 3", got)
	}

	// 23 - 3 = 20 left. Both minimums (5 + 5) come off first, leaving a surplus of
	// 10 split 10:30 -> 3 and 7. So 5+3 = 8 and 5+7 = 12.
	if got := flexA.GetSize().Height; got != 8 {
		t.Errorf("flexA height = %d, want 8", got)
	}

	if got := flexB.GetSize().Height; got != 12 {
		t.Errorf("flexB height = %d, want 12", got)
	}

	if got := sum(allocatedHeights(fixed, flexA, flexB)); got != 23 {
		t.Errorf("total allocated height = %d, want 23", got)
	}
}

// Rounding must never lose or invent rows: the last flexible element absorbs
// the remainder so the allocation always adds up to the available height.
func TestResizeFlexibleElementsFillsHeightExactly(t *testing.T) {
	for available := 10; available <= 60; available++ {
		fixed := newFlexStub(orvyn.NewSize(10, 3), orvyn.NewSize(10, 3))
		flexA := newFlexStub(orvyn.NewSize(10, 5), orvyn.NewSize(10, 7))
		flexB := newFlexStub(orvyn.NewSize(10, 5), orvyn.NewSize(10, 11))
		flexC := newFlexStub(orvyn.NewSize(10, 5), orvyn.NewSize(10, 13))

		resizeFlexibleElements(40, available, fixed, flexA, flexB, flexC)

		if got := sum(allocatedHeights(fixed, flexA, flexB, flexC)); got != available {
			t.Errorf("available %d: total allocated = %d, want %d", available, got, available)
		}
	}
}

// A flexible element allocated less than its minimal height renders at its
// minimum anyway, so the layout overflows by the difference. Never allocate below
// the minimum while there is room for it.
func TestResizeFlexibleElementsNeverAllocatesBelowMinimum(t *testing.T) {
	// Small preferred height next to a huge one: a pure preferred-height ratio
	// starves this element well under its minimum.
	small := newFlexStub(orvyn.NewSize(10, 10), orvyn.NewSize(10, 12))
	hungryA := newFlexStub(orvyn.NewSize(10, 5), orvyn.NewSize(10, 35))
	hungryB := newFlexStub(orvyn.NewSize(10, 5), orvyn.NewSize(10, 35))

	resizeFlexibleElements(40, 38, small, hungryA, hungryB)

	for name, e := range map[string]*flexStub{"small": small, "hungryA": hungryA, "hungryB": hungryB} {
		if got, want := e.GetSize().Height, e.GetMinSize().Height; got < want {
			t.Errorf("%s height = %d, below its minimum of %d", name, got, want)
		}
	}

	if got := sum(allocatedHeights(small, hungryA, hungryB)); got != 38 {
		t.Errorf("total allocated = %d, want 38", got)
	}
}

// When the minimums do not even fit, they are shared proportionally rather than
// letting the first elements take everything and the last get nothing.
func TestResizeFlexibleElementsSharesWhenMinimumsDoNotFit(t *testing.T) {
	a := newFlexStub(orvyn.NewSize(10, 10), orvyn.NewSize(10, 20))
	b := newFlexStub(orvyn.NewSize(10, 10), orvyn.NewSize(10, 20))

	resizeFlexibleElements(40, 10, a, b)

	if got := a.GetSize().Height; got != 5 {
		t.Errorf("a height = %d, want 5", got)
	}

	if got := b.GetSize().Height; got != 5 {
		t.Errorf("b height = %d, want 5", got)
	}
}

func TestResizeFlexibleElementsNeverAllocatesNegativeHeight(t *testing.T) {
	fixed := newFlexStub(orvyn.NewSize(10, 20), orvyn.NewSize(10, 20))
	flex := newFlexStub(orvyn.NewSize(10, 5), orvyn.NewSize(10, 30))

	resizeFlexibleElements(40, 5, fixed, flex)

	if got := flex.GetSize().Height; got < 0 {
		t.Errorf("flexible height = %d, want >= 0", got)
	}
}

func dashboardLikeLayout() (*DefinedWidthVerticalLayout, []orvyn.Renderable) {
	elements := []orvyn.Renderable{
		newFlexStub(orvyn.NewSize(10, 3), orvyn.NewSize(10, 3)),   // runningTask
		newFlexStub(orvyn.NewSize(30, 7), orvyn.NewSize(45, 7)),   // characterInfo
		newFlexStub(orvyn.NewSize(30, 10), orvyn.NewSize(45, 12)), // locationInfo
		newFlexStub(orvyn.NewSize(10, 5), orvyn.NewSize(30, 35)),  // logEvent
		newFlexStub(orvyn.NewSize(20, 5), orvyn.NewSize(60, 35)),  // social pile
		newFlexStub(orvyn.NewSize(1, 1), orvyn.NewSize(1, 1)),     // statusMessage
		newFlexStub(orvyn.NewSize(1, 1), orvyn.NewSize(1, 1)),     // help
	}

	return NewDefinedWidthVerticalLayout(35, 120, orvyn.NewSize(10, 10), elements...), elements
}

// The layout must hand out exactly the height it was given minus the vertical
// margin - no more (content gets cut off), no less (dead space). The floor is the
// combined height of the fixed elements, which cannot shrink.
func TestDefinedWidthVerticalFillsHeightExactly(t *testing.T) {
	const fixedTotal = 3 + 7 + 1 + 1

	l, elements := dashboardLikeLayout()

	for _, term := range []orvyn.Size{
		orvyn.NewSize(200, 60),
		orvyn.NewSize(120, 40),
		orvyn.NewSize(80, 24),
		orvyn.NewSize(40, 20),
	} {
		l.Resize(term)
		l.Render()

		want := max(term.Height-10, fixedTotal)

		if got := sum(allocatedHeights(elements...)); got != want {
			t.Errorf("terminal %dx%d: allocated height = %d, want %d",
				term.Width, term.Height, got, want)
		}
	}
}

// Shrinking the terminal must shrink the render, never leave it overflowing.
func TestDefinedWidthVerticalNeverOverflows(t *testing.T) {
	l, _ := dashboardLikeLayout()

	for width := 20; width <= 200; width += 7 {
		for height := 12; height <= 60; height += 5 {
			l.Resize(orvyn.NewSize(width, height))

			out := l.Render()
			lines := strings.Split(out, "\n")

			if len(lines) > height {
				t.Errorf("terminal %dx%d: rendered %d lines, available %d",
					width, height, len(lines), height)
			}

			for _, line := range lines {
				if got := len([]rune(line)); got > width {
					t.Errorf("terminal %dx%d: rendered line of %d columns, available %d",
						width, height, got, width)
					break
				}
			}
		}
	}
}

// Width policy: preferred width minus margin when there is room, clamped down to
// what the terminal actually offers when there is not.
func TestDefinedWidthVerticalWidthClamp(t *testing.T) {
	tests := []struct {
		layoutWidth int
		want        int
	}{
		{200, 110}, // capped at preferred (120) - margin (10)
		{120, 110},
		{80, 70}, // available - margin
		{40, 30}, // available - margin, still above min
		{30, 25}, // floored at min width (35) - margin
		{20, 20}, // floor is wider than the terminal: take what we have
	}

	for _, tc := range tests {
		l, elements := dashboardLikeLayout()

		l.Resize(orvyn.NewSize(tc.layoutWidth, 40))
		l.Render()

		if got := elements[0].GetSize().Width; got != tc.want {
			t.Errorf("layout width %d: element width = %d, want %d",
				tc.layoutWidth, got, tc.want)
		}
	}
}

// Inactive elements are not rendered, so they must not reserve any height.
func TestDefinedWidthVerticalIgnoresInactiveElements(t *testing.T) {
	fixed := newFlexStub(orvyn.NewSize(10, 4), orvyn.NewSize(10, 4))
	hidden := newFlexStub(orvyn.NewSize(10, 6), orvyn.NewSize(10, 6))
	flex := newFlexStub(orvyn.NewSize(10, 5), orvyn.NewSize(10, 30))

	l := NewDefinedWidthVerticalLayout(10, 40, orvyn.NewSize(0, 0), fixed, hidden, flex)

	hidden.SetActive(false)

	l.Resize(orvyn.NewSize(40, 30))
	l.Render()

	if got := flex.GetSize().Height; got != 26 {
		t.Errorf("flexible height = %d, want 26 (hidden element must not reserve height)", got)
	}
}

// Widgets change their own min/preferred size at runtime (a running task, a status
// message that gets text). The fixed/flexible split must follow, not stay frozen
// at construction time.
func TestDefinedWidthVerticalReclassifiesOnRender(t *testing.T) {
	morphing := newFlexStub(orvyn.NewSize(10, 3), orvyn.NewSize(10, 3))
	flex := newFlexStub(orvyn.NewSize(10, 5), orvyn.NewSize(10, 30))

	l := NewDefinedWidthVerticalLayout(10, 40, orvyn.NewSize(0, 0), morphing, flex)

	l.Resize(orvyn.NewSize(40, 30))
	l.Render()

	if got := flex.GetSize().Height; got != 27 {
		t.Errorf("flexible height = %d, want 27", got)
	}

	// The widget grew: it now reports a taller fixed height.
	morphing.min = orvyn.NewSize(10, 7)
	morphing.pref = orvyn.NewSize(10, 7)

	l.Render()

	if got := morphing.GetSize().Height; got != 7 {
		t.Errorf("morphing element height = %d, want 7", got)
	}

	if got := flex.GetSize().Height; got != 23 {
		t.Errorf("flexible height = %d, want 23 after the fixed element grew", got)
	}
}
