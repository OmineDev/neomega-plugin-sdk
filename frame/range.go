package frame

// Range is an inclusive minimum and maximum Y coordinate.
type Range [2]int

// Min returns the minimum coordinate.
func (r Range) Min() int {
	return r[0]
}

// Max returns the maximum coordinate.
func (r Range) Max() int {
	return r[1]
}

// Height returns Max minus Min; it does not add one for an inclusive range.
func (r Range) Height() int {
	return r[1] - r[0]
}
