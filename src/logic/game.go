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
	shipsMarkedForReturn map[int]bool    // keep track of ships returning to a dock
	dropOffs             []*hlt.Position // keep track of drop offs
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
	// If we have enough halite check if we should convert to a drop off
	if gm.game.Me.Halite > (dropCost * 2) {
		// this should never be true because we never let the ship get full
		// plus the convert logic was pulled out into a ConvertAI object
		if gm.shouldConvert(ship) {
			gm.dropOffs = append(gm.dropOffs, ship.E.Pos)
			return Convert
		}
	}
	// check our distance compared to how long it will take to get back to decide if we should return to a dock
	if _, d := gm.closestDropoff(ship.E.Pos); (d + 3) >= (maxTurn - gm.game.TurnNumber) {
		return Return
	}
	// if there is not enough halite to move and we are not on a drop off stay put
	if math.Ceil(float64(currentCell.Halite)*(1.0/moveCost)) > float64(ship.Halite) && !gm.onDropOff(ship.E.Pos) {
		return Stay
	}
	// check if ship is marked for return
	if t, ok := gm.shipsMarkedForReturn[ship.E.ID()]; ok && t {
		// if the ship has lost too much halite before hitting dock, forget about returning to dock
		// this logic masked the problem with the method hasShipReturned
		if (float64(ship.Halite)/float64(maxHalite)) < 0.4 && (maxTurn-gm.game.TurnNumber) > 50 {
			delete(gm.shipsMarkedForReturn, ship.E.ID())
			if currentCell.Halite < 10 {
				return Collect
			}
			return Stay
		}
		// if the ship has returned to a drop off it needs to move on
		// hasShipReturned implementation only checks for starting dock and not all docks
		// but the logic to create extra docks made it to where this rarely occurred, but still a bug
		if gm.hasShipReturned(currentCell, ship) {
			return Collect
		}
		return Return
	}
	// if the ship was full it wouldn't move
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
	// this should actually be changed to check for all dock positions
	if currentCell.Pos.Equals(gm.game.Me.Shipyard.E.Pos) {
		gm.shipsMarkedForReturn[ship.E.ID()] = false
		return true
	}
	return false
}

// this should never return true because we have logic to not let a ship become full
// we also have ConvertAI object to handle the logic for conversion
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
