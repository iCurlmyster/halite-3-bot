package logic

import (
	"fmt"
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
	cell := move.findMostHaliteInWindow(ship.E.Pos, 4+int(math.Floor(float64(move.gameAI.game.TurnNumber)/100.0)))

	return ship.Move(move.lazyGreedySearch(cell.Pos, ship.E.Pos, 4))
	// if d, found := move.findDirectionToCell(cell, ship.E.Pos); found {
	// 	return ship.Move(d)
	// }
	// return move.randomAvailableDirection(ship)
}

// function not used anymore because we switched to lazyGreedySearch
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
	// if the number of turns is less than the amount of ships
	// disregard saftey checks in lazyGreedySearch and just navigate to destination (should be dock)
	// even if they need to crash into other ships
	if (maxTurns - move.gameAI.game.TurnNumber) <= len(move.Me.Ships) {
		dirs := move.Map.GetUnsafeMoves(ship.E.Pos, dropoff)
		fDir := move.determineBestDirectionOutOfTwo(dirs, dropoff)
		if fDir != nil {
			return ship.Move(fDir)
		}
		return ship.StayStill()
	}
	return ship.Move(move.lazyGreedySearch(dropoff, ship.E.Pos, 4))
	// dir := move.Map.NaiveNavigate(ship, dropoff)
	// nextPos := helper.NormalizedDirectionalOffset(ship.E.Pos, move.Map, dir)
	// if !move.IsFutureClaimed(nextPos) {
	// 	move.MarkFuturePos(nextPos)
	// 	return ship.Move(dir)
	// }
	// move.MarkFuturePos(ship.E.Pos)
	// return ship.StayStill()
}

// create grid and keep searching out to a certain depth for the cell with the most halite and that is close
func (move *MoveAI) findMostHaliteInWindow(pos *hlt.Position, n int) *hlt.MapCell {
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
	return answer
}

// This method was replaced with the lazyGreedySearch
func (move *MoveAI) findDirectionToCell(answer *hlt.MapCell, pos *hlt.Position) (*hlt.Direction, bool) {
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
	return move.AvailableDirectionsForPos(cell.Pos)
}

// AvailableDirectionsForPos - Returns array of immediately available neighboring positions
func (move *MoveAI) AvailableDirectionsForPos(pos *hlt.Position) []*hlt.Direction {
	up := helper.NormalizedDirectionalOffset(pos, move.Map, hlt.North())
	down := helper.NormalizedDirectionalOffset(pos, move.Map, hlt.South())
	left := helper.NormalizedDirectionalOffset(pos, move.Map, hlt.West())
	right := helper.NormalizedDirectionalOffset(pos, move.Map, hlt.East())
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

// AvailablePositionsForPos -
func (move *MoveAI) AvailablePositionsForPos(pos *hlt.Position) []*hlt.Position {
	up := helper.NormalizedDirectionalOffset(pos, move.Map, hlt.North())
	down := helper.NormalizedDirectionalOffset(pos, move.Map, hlt.South())
	left := helper.NormalizedDirectionalOffset(pos, move.Map, hlt.West())
	right := helper.NormalizedDirectionalOffset(pos, move.Map, hlt.East())
	dirs := make([]*hlt.Position, 0)
	if !move.IsPosClaimed(up) {
		dirs = append(dirs, up)
	}
	if !move.IsPosClaimed(down) {
		dirs = append(dirs, down)
	}
	if !move.IsPosClaimed(left) {
		dirs = append(dirs, left)
	}
	if !move.IsPosClaimed(right) {
		dirs = append(dirs, right)
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

// perform a breadth for search up to a certain depth or if we find the destination
func (move *MoveAI) lazyGreedySearch(target, src *hlt.Position, depth int) *hlt.Direction {
	history := make(map[*hlt.Direction][]*hlt.Position)
	claimed := make(map[string]bool)
	parents := move.AvailableDirectionsForPos(src)
	for _, d := range parents {
		posDir := helper.NormalizedDirectionalOffset(src, move.Map, d)
		if posDir.Equals(target) {
			move.MarkFuturePos(posDir)
			return d
		}
		history[d] = []*hlt.Position{posDir}
	}
	for i := 0; i < depth; i++ {
		for _, d := range parents {
			his := history[d]
			curSrc := his[len(his)-1]
			avDirs := move.AvailablePositionsForPos(curSrc)
			bestDis := 10000000
			var bestPos *hlt.Position
			for _, ad := range avDirs {
				if _, ok := claimed[fmt.Sprintf("%s%s", d.String(), ad.String())]; !ok {
					if ad.Equals(target) {
						move.MarkFuturePos(his[0])
						return d
					}
					tmpDis := move.Map.CalculateDistance(ad, target)
					if tmpDis < bestDis {
						bestPos = ad
						bestDis = tmpDis
					}
				}
			}
			if bestPos != nil {
				history[d] = append(his, bestPos)
				claimed[fmt.Sprintf("%s%s", d.String(), bestPos.String())] = true
			}
		}
	}

	var bestDir *hlt.Direction
	bestDis := 100000000
	var bestPos *hlt.Position
	for k, v := range history {
		if len(v) == 0 {
			continue
		}
		last := v[len(v)-1]
		tmpDis := move.Map.CalculateDistance(last, target)
		if tmpDis < bestDis {
			bestDis = tmpDis
			bestDir = k
			bestPos = v[0]
		}
	}
	if bestDir != nil {
		move.MarkFuturePos(bestPos)
		return bestDir
	}
	return hlt.Still()
}
