package logic

import (
	"helper"
	"hlt"
	"math/rand"
)

// MoveAI - Object to store and compute logic for the game
type MoveAI struct {
	gameAI    *GameAI
	FuturePos map[string]*hlt.MapCell
	Map       *hlt.GameMap
	Me        *hlt.Player
}

// NewMoveAI - Generates a new MoveAI object to use
func NewMoveAI(gm *GameAI, gMap *hlt.GameMap, p *hlt.Player) *MoveAI {
	return &MoveAI{
		gameAI:    gm,
		FuturePos: make(map[string]*hlt.MapCell),
		Map:       gMap,
		Me:        p,
	}
}

// Move - Returns the command on how the ship should or should not move
func (move *MoveAI) Move(ship *hlt.Ship) hlt.Command {
	switch move.gameAI.ShipLogic(ship) {
	case Collect:
		return move.randomAvailableDirection(ship)
	case Return:
		dir := move.Map.NaiveNavigate(ship, move.Me.Shipyard.E.Pos)
		nextPos := helper.NormalizedDirectionalOffset(ship.E.Pos, move.Map, dir)
		if !move.IsPosClaimed(nextPos) {
			move.MarkFuturePos(nextPos)
			return ship.Move(dir)
		}
	case Convert:
		return ship.MakeDropoff()
	case Stay:
		break
	}
	move.MarkFuturePos(ship.E.Pos)
	return ship.StayStill()
}

// MarkFuturePos - Mark future positions of ship movements
func (move *MoveAI) MarkFuturePos(pos *hlt.Position) {
	cell := move.Map.AtPosition(pos)
	move.FuturePos[pos.String()] = cell
}

// IsPosClaimed - check to see if position has been claimed as a future spot
func (move *MoveAI) IsPosClaimed(pos *hlt.Position) bool {
	_, ok := move.FuturePos[pos.String()]
	return ok
}

func (move *MoveAI) randomAvailableDirection(ship *hlt.Ship) hlt.Command {
	dirs := helper.AvailableDirectionsForEntity(move.Map, ship.E)
	choice, nextPos := exhaustOptions(move, ship.E.Pos, dirs)
	move.MarkFuturePos(nextPos)
	return ship.Move(choice)
}

func exhaustOptions(move *MoveAI, pos *hlt.Position, dirs []*hlt.Direction) (*hlt.Direction, *hlt.Position) {
	if len(dirs) > 0 {
		// pick a random starting direction
		N := rand.Intn(len(dirs))
		choice := dirs[N]
		nextPos := helper.NormalizedDirectionalOffset(pos, move.Map, choice)
		// if not claimed use it.
		if !move.IsPosClaimed(nextPos) {
			return choice, nextPos
		}
		// if the spot was claimed exhaust other options
		for i := 0; i < len(dirs); i++ {
			if i == N {
				continue
			}
			choice = dirs[N]
			nextPos = helper.NormalizedDirectionalOffset(pos, move.Map, choice)
			if !move.IsPosClaimed(nextPos) {
				return choice, nextPos
			}
		}
	}
	// if there is no where to go, stay still
	return hlt.Still(), pos
}
