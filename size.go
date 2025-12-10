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

// DivideSizeFull is a helper function that allow to divide a size in 2. The first size returned will have the compensentation to avoid float result.
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
