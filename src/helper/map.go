package helper

import (
	"hlt"
)

// AvailableDirectionsForEntity - Returns array of immediately available neighboring positions
func AvailableDirectionsForEntity(gMap *hlt.GameMap, e *hlt.Entity) []*hlt.Direction {
	cell := gMap.AtEntity(e)
	up := NormalizedDirectionalOffset(cell.Pos, gMap, hlt.North())
	down := NormalizedDirectionalOffset(cell.Pos, gMap, hlt.South())
	left := NormalizedDirectionalOffset(cell.Pos, gMap, hlt.West())
	right := NormalizedDirectionalOffset(cell.Pos, gMap, hlt.East())
	dirs := make([]*hlt.Direction, 0)
	if PositionOpen(gMap, up) {
		dirs = append(dirs, hlt.North())
	}
	if PositionOpen(gMap, down) {
		dirs = append(dirs, hlt.South())
	}
	if PositionOpen(gMap, left) {
		dirs = append(dirs, hlt.West())
	}
	if PositionOpen(gMap, right) {
		dirs = append(dirs, hlt.East())
	}
	return dirs
}

// PositionOpen - Checks to see if the position is currently occupied
func PositionOpen(gMap *hlt.GameMap, p *hlt.Position) bool {
	c := gMap.AtPosition(p)
	return !c.IsOccupied()
}

// NormalizedDirectionalOffset - Get normalized position of direction offset
func NormalizedDirectionalOffset(pos *hlt.Position, gMap *hlt.GameMap, d *hlt.Direction) *hlt.Position {
	off, err := pos.DirectionalOffset(d)
	// this should never hit, If I am using it correctly
	if err != nil {
		panic(err)
	}
	return gMap.Normalize(off)
}
