package logic

import (
	"hlt"
	"hlt/gameconfig"
	"math"
)

const (
	dropoffDistancePoint = 13
)

var (
	dropoffDistanceThreshold = 2
)

// ConvertAI - Object to handle conversion of ships to drop offs
type ConvertAI struct {
	game           *GameAI
	CurrentDropoff *hlt.Ship
}

// NewConvertAI - Generates a new ConvertAI object
func NewConvertAI(g *GameAI) *ConvertAI {
	return &ConvertAI{
		game: g,
	}
}

// DeterminePossibleDropOff - Determines whether a ship can become a drop off or not
func (ca *ConvertAI) DeterminePossibleDropOff(ships map[int]*hlt.Ship) hlt.Command {
	dropCost, _ := ca.game.config.GetInt(gameconfig.DropoffCost)
	maxTurn, _ := ca.game.config.GetInt(gameconfig.MaxTurns)
	if ca.game.game.Me.Halite < (dropCost*(2*len(ca.game.dropOffs))) || (maxTurn-ca.game.game.TurnNumber) <= 230 {
		return nil
	}
	var optShip *hlt.Ship
	optDis := 0
	for i := range ships {
		var s = ships[i]
		if d, ok := ca.correctDistance(optDis, s.E.Pos, ca.game.dropOffs); ok {
			optShip = s
			optDis = d
		}
	}
	if optShip != nil && int(math.Abs(float64(dropoffDistancePoint-optDis))) <= dropoffDistanceThreshold {
		if ca.game.onDropOff(optShip.E.Pos) {
			return nil
		}
		ca.CurrentDropoff = optShip
		ca.game.dropOffs = append(ca.game.dropOffs, optShip.E.Pos)
		return optShip.MakeDropoff()
	}
	return nil
}

// IsCurrentDropoff - Checks to see if the ship is the converted ship for this turn
func (ca *ConvertAI) IsCurrentDropoff(s *hlt.Ship) bool {
	if ca.CurrentDropoff != nil && s != nil {
		return ca.CurrentDropoff.E.ID() == s.E.ID()
	}
	return false
}

func (ca *ConvertAI) correctDistance(optDis int, shipPos *hlt.Position, dos []*hlt.Position) (int, bool) {
	var curPos *hlt.Position
	curDis := optDis
	for _, d := range dos {
		tmpDis := ca.game.game.Map.CalculateDistance(shipPos, d)
		if (dropoffDistancePoint - tmpDis) > dropoffDistanceThreshold {
			return 0, false
		} else if curDis == 0 || tmpDis < curDis {
			curDis = tmpDis
			curPos = d
		}
	}
	if curPos != nil {
		return curDis, true
	}
	return 0, false
}
