package logic

import (
	"helper"
	"hlt"
	"math/rand"
)

// GameAI - Object to store and compute logic for the game
type GameAI struct {
	FuturePos map[string]*hlt.MapCell
	Map       *hlt.GameMap
	Me        *hlt.Player
}

// NewGameAI - Generates a new GameAI object to use
func NewGameAI(gMap *hlt.GameMap, p *hlt.Player) *GameAI {
	return &GameAI{
		FuturePos: make(map[string]*hlt.MapCell),
		Map:       gMap,
		Me:        p,
	}
}

// Move - Returns the command on how the ship should or should not move
func (gm *GameAI) Move(ship *hlt.Ship, maxHalite int) hlt.Command {
	var currentCell = gm.Map.AtEntity(ship.E)
	if currentCell.Halite < (maxHalite/10) || ship.IsFull() {
		return gm.randomAvailableDirection(ship)
	}
	gm.MarkFuturePos(ship.E.Pos)
	return ship.Move(hlt.Still())
}

// MarkFuturePos - Mark future positions of ship movements
func (gm *GameAI) MarkFuturePos(pos *hlt.Position) {
	cell := gm.Map.AtPosition(pos)
	gm.FuturePos[pos.String()] = cell
}

// IsPosClaimed - check to see if position has been claimed as a future spot
func (gm *GameAI) IsPosClaimed(pos *hlt.Position) bool {
	_, ok := gm.FuturePos[pos.String()]
	return ok
}

func (gm *GameAI) randomAvailableDirection(ship *hlt.Ship) hlt.Command {
	dirs := helper.AvailableDirectionsForEntity(gm.Map, ship.E)
	choice, nextPos := exhaustOptions(gm, ship.E.Pos, dirs)
	gm.MarkFuturePos(nextPos)
	return ship.Move(choice)
}

func exhaustOptions(gm *GameAI, pos *hlt.Position, dirs []*hlt.Direction) (*hlt.Direction, *hlt.Position) {
	if len(dirs) > 0 {
		// pick a random starting direction
		N := rand.Intn(len(dirs))
		choice := dirs[N]
		nextPos := helper.NormalizedDirectionalOffset(pos, gm.Map, choice)
		// if not claimed use it.
		if !gm.IsPosClaimed(nextPos) {
			return choice, nextPos
		}
		// if the spot was claimed exhaust other options
		for i := 0; i < len(dirs); i++ {
			if i == N {
				continue
			}
			choice = dirs[N]
			nextPos = helper.NormalizedDirectionalOffset(pos, gm.Map, choice)
			if !gm.IsPosClaimed(nextPos) {
				return choice, nextPos
			}
		}
	}
	// if there is no where to go, stay still
	return hlt.Still(), pos
}
