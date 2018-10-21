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
		return gm.randomAvailableDirection(gMap, ship)
	}
	gm.MarkFuturePos(gMap, ship.E.Pos)
	return ship.Move(hlt.Still())
}

// MarkFuturePos - Mark future positions of ship movements
func (gm *GameAI) MarkFuturePos(gMap *hlt.GameMap, pos *hlt.Position) {
	cell := gMap.AtPosition(pos)
	gm.FuturePos[pos.String()] = cell
}

// IsPosClaimed - check to see if position has been claimed as a future spot
func (gm *GameAI) IsPosClaimed(pos *hlt.Position) bool {
	_, ok := gm.FuturePos[pos.String()]
	return ok
}

func (gm *GameAI) randomAvailableDirection(gMap *hlt.GameMap, ship *hlt.Ship) hlt.Command {
	dirs := helper.AvailableDirectionsForEntity(gMap, ship.E)
	choice, nextPos := exhaustOptions(gm, gMap, ship.E.Pos, dirs)
	gm.MarkFuturePos(gMap, nextPos)
	return ship.Move(choice)
}

func exhaustOptions(gm *GameAI, gMap *hlt.GameMap, pos *hlt.Position, dirs []*hlt.Direction) (*hlt.Direction, *hlt.Position) {
	if len(dirs) > 0 {
		N := rand.Intn(len(dirs))
		choice := dirs[N]
		nextPos := helper.NormalizedDirectionalOffset(pos, gMap, choice)
		if !gm.IsPosClaimed(nextPos) {
			return choice, nextPos
		}
		for i := 0; i < len(dirs); i++ {
			if i == N {
				continue
			}
			choice = dirs[N]
			nextPos = helper.NormalizedDirectionalOffset(pos, gMap, choice)
			if !gm.IsPosClaimed(nextPos) {
				return choice, nextPos
			}
		}
	}
	return hlt.Still(), pos
}
