package main

import (
	"fmt"
	"hlt"
	"hlt/gameconfig"
	"hlt/log"
	"logic"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func gracefulExit(logger *log.FileLogger) {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v", sig)
		fmt.Println("Wait for 2 second to finish processing")
		time.Sleep(2 * time.Second)
		logger.Close()
		os.Exit(0)
	}()
}

func main() {
	args := os.Args
	var seed = time.Now().UnixNano() % int64(os.Getpid())
	if len(args) > 1 {
		seed, _ = strconv.ParseInt(args[0], 10, 64)
	}
	rand.Seed(seed)

	var game = hlt.NewGame()
	// At this point "game" variable is populated with initial map data.
	// This is a good place to do computationally expensive start-up pre-processing.
	// As soon as you call "ready" function below, the 2 second per turn timer will start.

	// TODO scan board with a window to build a hueristic of most desirable locations on map. Set in GameAI
	// determine window by taking  width/8  to yield window size. example 32x32 map -> 32/8 yields window of size 4

	var config = gameconfig.GetInstance()
	// Setup GameAI to persist data between frames
	gameAI := logic.NewGameAI(game, config)
	maxShipCount := 8

	fileLogger := log.NewFileLogger(game.Me.ID)
	var logger = fileLogger.Logger
	logger.Printf("Successfully created bot! My Player ID is %d. Bot rng seed is %d.", game.Me.ID, seed)
	gracefulExit(fileLogger)
	game.Ready("jm")
	maxTurn, _ := config.GetInt(gameconfig.MaxTurns)
	for {
		game.UpdateFrame()
		var me = game.Me
		var gameMap = game.Map
		var ships = me.Ships
		var commands = []hlt.Command{}
		var moveAI = logic.NewMoveAI(gameAI, gameMap, me)
		var convertAI = logic.NewConvertAI(gameAI)
		if com := convertAI.DeterminePossibleDropOff(ships); com != nil {
			commands = append(commands, com)
		}
		for i := range ships {
			var ship = ships[i]
			if convertAI.IsCurrentDropoff(ship) {
				continue
			}
			commands = append(commands, moveAI.Move(ship))
		}
		var shipCost, _ = config.GetInt(gameconfig.ShipCost)
		if len(ships) < maxShipCount && me.Halite >= (3*shipCost) && !gameMap.AtEntity(me.Shipyard.E).IsOccupied() && (maxTurn-game.TurnNumber) <= 100 {
			commands = append(commands, hlt.SpawnShip{})
			if (len(ships)+1) >= maxShipCount && maxShipCount > 6 {
				maxShipCount--
			}
		}
		if math.Mod(float64(game.TurnNumber), 100.0) == 0 {
			maxShipCount--
		}
		game.EndTurn(commands)
	}
}
