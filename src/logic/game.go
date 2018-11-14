package logic

import (
	"hlt"
	"hlt/gameconfig"
	"math"
)

// ShipDecision - Representation of what logic the ship should perform
type ShipDecision int

const (
	// Collect - Collect more halite
	Collect = ShipDecision(iota)
	// Return - Return to nearest drop off
	Return
	// Convert - Convert ship into new drop off point
	Convert
	// Stay - The ship needs to stay where it is
	Stay
)

// GameAI - Object to store/handle overall game logic
type GameAI struct {
	game                 *hlt.Game
	config               *gameconfig.Constants
	shipsMarkedForReturn map[int]bool
	dropOffs             []*hlt.Position
}

// NewGameAI - Generate a new GameAI object
func NewGameAI(g *hlt.Game, c *gameconfig.Constants) *GameAI {
	dos := make([]*hlt.Position, 0)
	dos = append(dos, g.Me.Shipyard.E.Pos)
	return &GameAI{
		game:                 g,
		config:               c,
		shipsMarkedForReturn: make(map[int]bool),
		dropOffs:             dos,
	}
}

// ShipLogic - Figure out what decision the ship should make next
func (gm *GameAI) ShipLogic(ship *hlt.Ship) ShipDecision {
	currentCell := gm.game.Map.AtEntity(ship.E)
	maxHalite, _ := gm.config.GetInt(gameconfig.MaxHalite)
	moveCost, _ := gm.config.GetDouble(gameconfig.MoveCostRatio)
	dropCost, _ := gm.config.GetInt(gameconfig.DropoffCost)
	maxTurn, _ := gm.config.GetInt(gameconfig.MaxTurns)
	if gm.game.Me.Halite > (dropCost * 2) {
		if gm.shouldConvert(ship) {
			gm.dropOffs = append(gm.dropOffs, ship.E.Pos)
			return Convert
		}
	}
	if _, d := gm.closestDropoff(ship.E.Pos); (d + 3) >= (maxTurn - gm.game.TurnNumber) {
		return Return
	}
	if math.Ceil(float64(currentCell.Halite)*(1.0/moveCost)) > float64(ship.Halite) && !gm.onDropOff(ship.E.Pos) {
		return Stay
	}
	if t, ok := gm.shipsMarkedForReturn[ship.E.ID()]; ok && t {
		// if the ship has returned to a drop off it needs to move on
		if (float64(ship.Halite)/float64(maxHalite)) < 0.4 && (maxTurn-gm.game.TurnNumber) > 50 {
			delete(gm.shipsMarkedForReturn, ship.E.ID())
			if currentCell.Halite < 10 {
				return Collect
			}
			return Stay
		}
		if gm.hasShipReturned(currentCell, ship) {
			return Collect
		}
		return Return
	}
	if (float64(ship.Halite) / float64(maxHalite)) > 0.9 {
		gm.shipsMarkedForReturn[ship.E.ID()] = true
		return Return
	}
	if currentCell.Halite < 10 {
		return Collect
	}
	return Stay
}

func (gm *GameAI) onDropOff(pos *hlt.Position) bool {
	for i := 0; i < len(gm.dropOffs); i++ {
		if pos.Equals(gm.dropOffs[i]) {
			return true
		}
	}
	return false
}

func (gm *GameAI) hasShipReturned(currentCell *hlt.MapCell, ship *hlt.Ship) bool {
	if currentCell.Pos.Equals(gm.game.Me.Shipyard.E.Pos) {
		gm.shipsMarkedForReturn[ship.E.ID()] = false
		return true
	}
	return false
}

func (gm *GameAI) shouldConvert(ship *hlt.Ship) bool {
	// TODO set up logic to determine conversion
	if ship.IsFull() {
		return true
	}
	return false
}

func (gm *GameAI) closestDropoff(pos *hlt.Position) (*hlt.Position, int) {
	var curPos *hlt.Position
	curDis := 0
	for _, d := range gm.dropOffs {
		tmpDis := gm.game.Map.CalculateDistance(pos, d)
		if curPos == nil {
			curPos = d
			curDis = tmpDis
		} else if tmpDis < curDis {
			curDis = tmpDis
			curPos = d
		}
	}
	return curPos, curDis
}
