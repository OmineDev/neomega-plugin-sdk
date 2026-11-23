package frame

import (
	"fmt"
)

// Dimension identifies the Overworld, Nether or End.
type Dimension int32

const (
	// DimensionOverworld is the Overworld dimension.
	DimensionOverworld = Dimension(0)
	// DimensionNether is the Nether dimension.
	DimensionNether = Dimension(1)
	// DimensionTheEnd is the End dimension.
	DimensionTheEnd = Dimension(2)
)

// String returns a standard dimension name, or dmN for an unrecognized dimension value.
func (d Dimension) String() string {
	switch d {
	case DimensionOverworld:
		return "overworld"
	case DimensionNether:
		return "nether"
	case DimensionTheEnd:
		return "the_end"
	default:
		return fmt.Sprintf("dm%v", int32(d))
	}
}

// RangeUpperInclude returns the inclusive block-height range for the dimension. Unknown values use the
// Overworld range.
func (d Dimension) RangeUpperInclude() Range {
	switch d {
	case DimensionOverworld:
		return OverWorldRangeUpperInclude
	case DimensionNether:
		return NetherRangeUpperInclude
	case DimensionTheEnd:
		return TheEndRangeUpperInclude
	default:
		return OverWorldRangeUpperInclude
	}
}

// RangeUpperExclude returns the exclusive upper-bound range for the dimension. Unknown values use the
// Overworld range.
func (d Dimension) RangeUpperExclude() Range {
	switch d {
	case DimensionOverworld:
		return OverWorldRangeUpperExclude
	case DimensionNether:
		return NetherRangeUpperExclude
	case DimensionTheEnd:
		return TheEndRangeUpperExclude
	default:
		return OverWorldRangeUpperExclude
	}
}

// OverWorldRangeUpperInclude is the inclusive Overworld block-height range.
var OverWorldRangeUpperInclude = Range{-64, 319}

// NetherRangeUpperInclude is the inclusive Nether block-height range.
var NetherRangeUpperInclude = Range{0, 127}

// TheEndRangeUpperInclude is the inclusive End block-height range.
var TheEndRangeUpperInclude = Range{0, 255}

// OverWorldRangeUpperExclude is the Overworld block-height range with an exclusive upper bound.
var OverWorldRangeUpperExclude = Range{-64, 320}

// NetherRangeUpperExclude is the Nether block-height range with an exclusive upper bound.
var NetherRangeUpperExclude = Range{0, 128}

// TheEndRangeUpperExclude is the End block-height range with an exclusive upper bound.
var TheEndRangeUpperExclude = Range{0, 254}
