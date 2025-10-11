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

func DivideSizeFull(size int) (int, int) {
	var result int

	result = size / 2

	totalSize := result * 2

	if totalSize == size {
		return result, result
	}

	compensation := size - totalSize

	return result + compensation, result
}
