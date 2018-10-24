package helper

import (
	"hlt"
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
