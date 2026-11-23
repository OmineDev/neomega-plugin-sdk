package frame

import (
	"fmt"
)

// ChunkPos stores an X,Z coordinate in chunk units.
type ChunkPos [2]int32

// X returns the X chunk coordinate.
func (p ChunkPos) X() int32 {
	return p[0]
}

// Z returns the Z chunk coordinate.
func (p ChunkPos) Z() int32 {
	return p[1]
}

// String formats the chunk coordinate as (x, z).
func (p ChunkPos) String() string {
	return fmt.Sprintf("(%v, %v)", p[0], p[1])
}

// WorldRange is the inclusive block-height range used for the Overworld.
var WorldRange = Range{-64, 319}

// CubePos stores a block coordinate in X,Y,Z order. Coordinates are in blocks.
type CubePos [3]int

// SortStartAndEndPos returns two positions ordered independently on each axis.
func (c CubePos) SortStartAndEndPos(another CubePos) (start, end CubePos) {
	return SortStartAndEndPos(c, another)
}

// SortStartAndEndPos returns start and end with each axis ordered independently.
func SortStartAndEndPos(_start, _end CubePos) (start, end CubePos) {
	if _start.X() > _end.X() {
		start[0] = _end.X()
		end[0] = _start.X()
	} else {
		start[0] = _start.X()
		end[0] = _end.X()
	}
	if _start.Y() > _end.Y() {
		start[1] = _end.Y()
		end[1] = _start.Y()
	} else {
		start[1] = _start.Y()
		end[1] = _end.Y()
	}
	if _start.Z() > _end.Z() {
		start[2] = _end.Z()
		end[2] = _start.Z()
	} else {
		start[2] = _start.Z()
		end[2] = _end.Z()
	}
	return start, end
}

// CubeSize returns the absolute per-axis difference between this position and another.
func (c CubePos) CubeSize(another CubePos) CubePos {
	return CubeSize(c, another)
}

// CubeSize returns the absolute difference between start and end on each axis. It does not add one for
// inclusive bounds.
func CubeSize(start, end CubePos) CubePos {
	size := end.Sub(start)
	if size[0] < 0 {
		size[0] = -size[0]
	}
	if size[1] < 0 {
		size[1] = -size[1]
	}
	if size[2] < 0 {
		size[2] = -size[2]
	}
	return size
}

// OutOfYBounds reports whether Y lies outside the inclusive WorldRange.
func (p CubePos) OutOfYBounds() bool {
	y := p[1]
	return y > WorldRange[1] || y < WorldRange[0]
}

// String formats the block coordinate as (x,y,z).
func (p CubePos) String() string {
	return fmt.Sprintf("(%v,%v,%v)", p[0], p[1], p[2])
}

// Sub subtracts another position component by component.
func (p CubePos) Sub(po CubePos) (offset CubePos) {
	offset[0] = p[0] - po[0]
	offset[1] = p[1] - po[1]
	offset[2] = p[2] - po[2]
	return offset
}

// Add adds another position component by component.
func (p CubePos) Add(po CubePos) (offset CubePos) {
	offset[0] = p[0] + po[0]
	offset[1] = p[1] + po[1]
	offset[2] = p[2] + po[2]
	return offset
}

// X returns the X block coordinate.
func (p CubePos) X() int {
	return p[0]
}

// Y returns the Y block coordinate.
func (p CubePos) Y() int {
	return p[1]
}

// Z returns the Z block coordinate.
func (p CubePos) Z() int {
	return p[2]
}
