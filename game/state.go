// Package game contains the state and rules for TermZogs.
package game

import (
	"math/rand"
	"time"
)

const (
	// InitialEnergy and MaxEnergy define the starting and capacity values.
	InitialEnergy         = 200
	MaxEnergy             = 300
	StatusMessageLifetime = 6

	MazeViewMinCredits      = 30
	MazeViewCost            = 20
	CombatAnimationDuration = 2 * time.Second
)

// GameState identifies the current screen and whether play can continue.
type GameState uint8

const (
	// StateTitle waits for input; StatePlaying enables game actions.
	StateTitle GameState = iota
	StatePlaying
	StateVictory
	StateDefeat
	StateInstructions
)

type DefeatCause uint8

const (
	DefeatCauseNone DefeatCause = iota
	DefeatCauseEnergy
	DefeatCauseZog
)

// Player tracks the explorer's position, inventory, energy, and moves.
type Player struct {
	Position    Point
	Facing      Point
	HasSword    bool
	HasTreasure bool
	Energy      int
	MaxEnergy   int
	Moves       int
}

// Game combines maze data with player, combat, and guidance state.
type Game struct {
	Maze               *Maze
	Player             Player
	State              GameState
	DefeatCause        DefeatCause
	StatusMessage      string
	StatusMessageSteps int
	ZogsSlain          int
	PendingSwap        bool
	GuidancePath       []Point
	GuidanceTicks      int
	ExpandedView       bool
	combatStartedAt    time.Time
	combatPoint        Point
	rng                *rand.Rand
	seed               int64
}
