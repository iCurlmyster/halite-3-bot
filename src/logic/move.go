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
		return move.determinePath(ship)
	case Return:
		return move.navigateToDropOff(ship)
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
	c := move.Map.AtPosition(pos)
	return move.IsFutureClaimed(pos) || c.IsOccupied()
}

// IsFutureClaimed - check to see if a position has been claimed as a future position
func (move *MoveAI) IsFutureClaimed(pos *hlt.Position) bool {
	_, ok := move.FuturePos[pos.String()]
	return ok
}

func (move *MoveAI) determinePath(ship *hlt.Ship) hlt.Command {
	if c, found := move.findHaliteInWindow(ship.E.Pos, 4); found {
		return c
	}
	return move.randomAvailableDirection(ship)
}

func (move *MoveAI) randomAvailableDirection(ship *hlt.Ship) hlt.Command {
	dirs := move.AvailableDirectionsForEntity(ship.E)
	nextPos := ship.E.Pos
	choice := hlt.Still()
	if len(dirs) > 0 {
		choice = dirs[rand.Intn(len(dirs))]
		nextPos = helper.NormalizedDirectionalOffset(ship.E.Pos, move.Map, choice)
	}
	move.MarkFuturePos(nextPos)
	return ship.Move(choice)
}

func (move *MoveAI) navigateToDropOff(ship *hlt.Ship) hlt.Command {
	dir := move.Map.NaiveNavigate(ship, move.Me.Shipyard.E.Pos)
	nextPos := helper.NormalizedDirectionalOffset(ship.E.Pos, move.Map, dir)
	if !move.IsFutureClaimed(nextPos) {
		move.MarkFuturePos(nextPos)
		return ship.Move(dir)
	}
	move.MarkFuturePos(ship.E.Pos)
	return ship.StayStill()
}

func (move *MoveAI) findHaliteInWindow(pos *hlt.Position, n int) (hlt.Command, bool) {
	// TODO look for the best immediate path to head towards.
	return nil, false
}

// AvailableDirectionsForEntity - Returns array of immediately available neighboring positions
func (move *MoveAI) AvailableDirectionsForEntity(e *hlt.Entity) []*hlt.Direction {
	cell := move.Map.AtEntity(e)
	up := helper.NormalizedDirectionalOffset(cell.Pos, move.Map, hlt.North())
	down := helper.NormalizedDirectionalOffset(cell.Pos, move.Map, hlt.South())
	left := helper.NormalizedDirectionalOffset(cell.Pos, move.Map, hlt.West())
	right := helper.NormalizedDirectionalOffset(cell.Pos, move.Map, hlt.East())
	dirs := make([]*hlt.Direction, 0)
	if !move.IsPosClaimed(up) {
		dirs = append(dirs, hlt.North())
	}
	if !move.IsPosClaimed(down) {
		dirs = append(dirs, hlt.South())
	}
	if !move.IsPosClaimed(left) {
		dirs = append(dirs, hlt.West())
	}
	if !move.IsPosClaimed(right) {
		dirs = append(dirs, hlt.East())
	}
	return dirs
}
