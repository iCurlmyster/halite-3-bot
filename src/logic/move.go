package logic

import (
	"helper"
	"hlt"
	"math/rand"
)

// GameAI - Object to store and compute logic for the game
type GameAI struct {
	FuturePos map[string]*hlt.MapCell
}

// NewGameAI - Generates a new GameAI object to use
func NewGameAI() *GameAI {
	return &GameAI{
		FuturePos: make(map[string]*hlt.MapCell),
	}
}

// Move - Returns the command on how the ship should or should not move
func (gm *GameAI) Move(gMap *hlt.GameMap, ship *hlt.Ship, maxHalite int) hlt.Command {
	var currentCell = gMap.AtEntity(ship.E)
	if currentCell.Halite < (maxHalite/10) || ship.IsFull() {
		dirs := helper.AvailableDirectionsForEntity(gMap, ship.E)
		if len(dirs) > 0 {
			choice := dirs[rand.Intn(len(dirs))]
			gm.MarkFuturePos(gMap, helper.NormalizedDirectionalOffset(ship.E.Pos, gMap, choice))
			return ship.Move(dirs[rand.Intn(len(dirs))])
		}
	}
	gm.MarkFuturePos(gMap, ship.E.Pos)
	return ship.Move(hlt.Still())
}

// MarkFuturePos - Mark future positions of ship movements
func (gm *GameAI) MarkFuturePos(gMap *hlt.GameMap, pos *hlt.Position) {
	cell := gMap.AtPosition(pos)
	gm.FuturePos[pos.String()] = cell
}
