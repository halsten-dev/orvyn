package orvyn

// Size is a simple struct to represent a size.
type Size struct {
	Width  int
	Height int
}

// NewSize returns a new Size.
func NewSize(width, height int) Size {
	return Size{width, height}
}

func SameSize(s1, s2 Size) bool {
	if s1 == s2 {
		return true
	}

	return false
}
