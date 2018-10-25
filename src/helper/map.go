package helper

import (
	"hlt"
)

var (
	allDirs = []*hlt.Direction{hlt.East(), hlt.South(), hlt.West(), hlt.North()}
)

// NormalizedDirectionalOffset - Get normalized position of direction offset
func NormalizedDirectionalOffset(pos *hlt.Position, gMap *hlt.GameMap, d *hlt.Direction) *hlt.Position {
	off, err := pos.DirectionalOffset(d)
	// this should never hit, If I am using it correctly
	if err != nil {
		panic(err)
	}
	return gMap.Normalize(off)
}

// NormalizedGridOutlineOffset - Grabs the normalized positions in the outline of the grid depth specified
func NormalizedGridOutlineOffset(pos *hlt.Position, gMap *hlt.GameMap, depth int) []*hlt.Position {
	start := pos
	for i := 0; i < depth; i++ {
		start, _ = pos.DirectionalOffset(hlt.North())
	}
	for i := 0; i < depth; i++ {
		start, _ = pos.DirectionalOffset(hlt.West())
	}
	grid := make([]*hlt.Position, 0)
	row := (2 * depth)
	for i := 0; i < len(allDirs); i++ {
		dir := allDirs[i]
		for j := 0; j < row; j++ {
			start, _ = start.DirectionalOffset(dir)
			grid = append(grid, gMap.Normalize(start))
		}
	}
	return grid
}
