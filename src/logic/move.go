package logic

import (
	"helper"
	"hlt"
	"hlt/gameconfig"
	"math"
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
		move.gameAI.dropOffs = append(move.gameAI.dropOffs, ship.E.Pos)
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
	if d, found := move.findHaliteInWindow(ship.E.Pos, 4+int(math.Floor(float64(move.gameAI.game.TurnNumber)/100.0))); found {
		return ship.Move(d)
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
	maxTurns, _ := move.gameAI.config.GetInt(gameconfig.MaxTurns)
	var dropoff *hlt.Position
	dDis := 0
	for _, d := range move.gameAI.dropOffs {
		curDis := move.Map.CalculateDistance(ship.E.Pos, d)
		if dropoff == nil {
			dropoff = d
			dDis = curDis
		} else if curDis < dDis {
			dropoff = d
			dDis = curDis
		}
	}
	if (maxTurns - move.gameAI.game.TurnNumber) <= len(move.Me.Ships) {
		dirs := move.Map.GetUnsafeMoves(ship.E.Pos, dropoff)
		fDir := move.determineBestDirectionOutOfTwo(dirs, dropoff)
		if fDir != nil {
			return ship.Move(fDir)
		}
		return ship.StayStill()
	}
	dir := move.Map.NaiveNavigate(ship, dropoff)
	nextPos := helper.NormalizedDirectionalOffset(ship.E.Pos, move.Map, dir)
	if !move.IsFutureClaimed(nextPos) {
		move.MarkFuturePos(nextPos)
		return ship.Move(dir)
	}
	move.MarkFuturePos(ship.E.Pos)
	return ship.StayStill()
}

func (move *MoveAI) findHaliteInWindow(pos *hlt.Position, n int) (*hlt.Direction, bool) {
	// TODO look for the best immediate path to head towards.
	var answer *hlt.MapCell
	for i := 0; i < n; i++ {
		panels := helper.NormalizedGridOutlineOffset(pos, move.Map, i+1)
		for j := 0; j < len(panels); j++ {
			cell := move.Map.AtPosition(panels[j])
			if cell.Halite > 10 {
				if answer == nil {
					answer = cell
				} else if cell.Halite > answer.Halite {
					answer = cell
				} else if cell.Halite == answer.Halite {
					cp := move.Map.CalculateDistance(pos, cell.Pos)
					ap := move.Map.CalculateDistance(pos, answer.Pos)
					if cp < ap {
						answer = cell
					}
				}
			}
		}
	}
	if answer != nil {
		dirs := move.Map.GetUnsafeMoves(pos, answer.Pos)
		finalDir := move.determineBestDirectionOutOfTwo(dirs, pos)
		if finalDir != nil {
			nextPos := helper.NormalizedDirectionalOffset(pos, move.Map, finalDir)
			if !move.IsPosClaimed(nextPos) {
				move.MarkFuturePos(nextPos)
				return finalDir, true
			}
		}
	}
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

func (move *MoveAI) determineBestDirectionOutOfTwo(dirs []*hlt.Direction, pos *hlt.Position) *hlt.Direction {
	var finalDir *hlt.Direction
	dir1 := helper.NormalizedDirectionalOffset(pos, move.Map, dirs[0])
	dir2 := helper.NormalizedDirectionalOffset(pos, move.Map, dirs[1])
	dir1Dis := move.Map.CalculateDistance(pos, dir1)
	dir2Dis := move.Map.CalculateDistance(pos, dir2)
	if dir1Dis == 0 {
		finalDir = dirs[1]
	} else if dir2Dis == 0 {
		finalDir = dirs[0]
	} else if dir1Dis < dir2Dis {
		finalDir = dirs[0]
	} else {
		finalDir = dirs[1]
	}
	return finalDir
}
