// Package game contains player state and maze interaction rules.
package game

import (
	"math/rand"
	"time"
)

// NewGame generates a maze and places the player at its start.
func NewGame() *Game {
	return NewGameWithSeed(time.Now().UnixNano())
}

// NewGameWithSeed creates a playing game that can be reproduced from a seed.
func NewGameWithSeed(seed int64) *Game {
	current := newGame(seed)
	current.Start()
	return current
}

// NewTitleGame creates a game that waits at the title screen before generation.
func NewTitleGame() *Game {
	seed := time.Now().UnixNano()
	return &Game{State: StateTitle, seed: seed, rng: rand.New(rand.NewSource(seed))}
}

// newGame initializes seeded random state before a run begins.
func newGame(seed int64) *Game {
	return &Game{seed: seed, rng: rand.New(rand.NewSource(seed))}
}

// Start generates the maze and initializes the player for a new run.
func (current *Game) Start() {
	current.Maze = GenerateMazeWithSeed(current.seed)
	current.zogPatrols = createZogPatrols(current.Maze, current.seed)
	current.nextZogMoveAt = time.Now().Add(ZogMoveInterval)
	current.combatStartedAt = time.Time{}
	current.combatPoint = Point{}
	path, _ := current.Maze.FindPath(current.Maze.Start, current.Maze.Treasure)
	capacity := MaxEnergy
	if len(path)+2 > capacity {
		capacity = len(path) + 2
	}

	current.Player = Player{Position: current.Maze.Start, Facing: Point{Y: 1}, Energy: InitialEnergy, MaxEnergy: capacity}
	current.State = StatePlaying
	current.DefeatCause = DefeatCauseNone
	current.SetStatusMessage("")
	current.ZogsSlain = 0
	current.PendingSwap = false
	current.GuidancePath = nil
	current.GuidanceTicks = 0
	current.ExpandedView = false
}

// SetStatusMessage displays a message until six successful movement steps pass.
func (current *Game) SetStatusMessage(message string) {
	current.StatusMessage = message
	if message == "" {
		current.StatusMessageSteps = 0
		return
	}
	current.StatusMessageSteps = StatusMessageLifetime
}

func (current *Game) advanceStatusMessage() {
	if current.StatusMessage == "" || current.StatusMessageSteps <= 0 {
		return
	}
	current.StatusMessageSteps--
	if current.StatusMessageSteps == 0 {
		current.StatusMessage = ""
	}
}

// ToggleMazeView opens the expanded view if the player can afford its credit cost.
func (current *Game) ToggleMazeView() {
	if current.State != StatePlaying {
		return
	}
	if current.ExpandedView {
		current.ExpandedView = false
		current.SetStatusMessage("Returned to the local maze view.")
		return
	}
	if current.Player.Energy <= MazeViewMinCredits {
		current.SetStatusMessage("You don't have enough health credits to view the maze.")
		return
	}

	current.Player.Energy -= MazeViewCost
	current.ExpandedView = true
	current.SetStatusMessage("Expanded maze view opened (-20 credits).")
}

// Restart begins a fresh run with a new maze seed.
func (current *Game) Restart() {
	current.seed = time.Now().UnixNano()
	current.rng = rand.New(rand.NewSource(current.seed))
	current.Start()
}

// Move advances the player into a walkable tile and resolves its effects.
func (current *Game) Move(dx, dy int) {
	if current.State != StatePlaying || current.PendingSwap || current.CombatInProgress() {
		return
	}
	if absInt(dx)+absInt(dy) != 1 {
		return
	}
	current.Player.Facing = Point{X: dx, Y: dy}
	next := Point{X: current.Player.Position.X + dx, Y: current.Player.Position.Y + dy}
	if current.Maze.TileAt(next) == Prisoner {
		current.Interact()
		return
	}
	if !current.Maze.IsWalkable(next) {
		return
	}

	current.Player.Position = next
	current.Player.Moves++
	current.advanceStatusMessage()
	current.advanceGuidance()
	current.spendEnergy(1)
	if current.State == StateDefeat {
		return
	}

	switch current.Maze.TileAt(next) {
	case Sword:
		if current.Player.HasTreasure {
			current.PendingSwap = true
			current.SetStatusMessage("Drop Treasure to pick up Sword? [Y/N]")
		} else if !current.Player.HasSword {
			current.Player.HasSword = true
			current.Maze.Tiles[next.Y][next.X] = Floor
		}
	case Treasure:
		if !current.Player.HasSword {
			current.Player.HasTreasure = true
			current.Maze.Tiles[next.Y][next.X] = Floor
		}
	case Food:
		current.RestoreEnergy(50)
		current.Maze.Tiles[next.Y][next.X] = Floor
		current.SetStatusMessage("You found food and restored energy.")
	case Monster:
		current.beginCombat(next)
		return
	}

	if current.State == StatePlaying && current.Player.Position == current.Maze.Start && current.Player.HasTreasure {
		current.State = StateVictory
		current.SetStatusMessage("You escaped with the treasure!")
	}
}

// absInt returns the absolute value used to validate a one-cell move.
func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

// DropItem leaves the held sword or treasure on the player's current cell.
func (current *Game) DropItem() {
	if current.State != StatePlaying || current.PendingSwap {
		return
	}
	point := current.Player.Position
	switch {
	case current.Player.HasSword:
		current.Maze.Tiles[point.Y][point.X] = Sword
		current.Player.HasSword = false
		current.SetStatusMessage("You dropped the sword.")
	case current.Player.HasTreasure:
		current.Maze.Tiles[point.Y][point.X] = Treasure
		current.Player.HasTreasure = false
		current.SetStatusMessage("You dropped the treasure.")
	default:
		current.SetStatusMessage("You have nothing to drop.")
	}
}

// AnswerSwordSwap accepts or declines dropping treasure for a sword.
func (current *Game) AnswerSwordSwap(accept bool) {
	if !current.PendingSwap {
		return
	}
	current.PendingSwap = false
	if accept {
		point := current.Player.Position
		current.Maze.Tiles[point.Y][point.X] = Treasure
		current.Player.HasTreasure = false
		current.Player.HasSword = true
		current.SetStatusMessage("You dropped the treasure and took the sword.")
	} else {
		current.SetStatusMessage("You kept the treasure.")
	}
}

// Interact asks an adjacent prisoner for a route to the treasure or exit.
func (current *Game) Interact() {
	if current.State != StatePlaying || current.PendingSwap {
		return
	}
	if !current.adjacentToPrisoner() {
		current.SetStatusMessage("There is no prisoner nearby.")
		return
	}
	current.advanceGuidance()
	current.spendEnergy(15)
	if current.State == StateDefeat {
		return
	}

	destination := current.Maze.Treasure
	if current.Player.HasTreasure {
		destination = current.Maze.Start
	}
	path, found := current.Maze.FindPath(current.Player.Position, destination)
	if !found {
		current.SetStatusMessage("The prisoner cannot find a route.")
		return
	}
	current.GuidancePath = path
	current.GuidanceTicks = 15
	current.SetStatusMessage("The prisoner shows you the way.")
}

// adjacentToPrisoner reports whether a prisoner is in a cardinally adjacent cell.
func (current *Game) adjacentToPrisoner() bool {
	for _, direction := range [...]Point{{X: 0, Y: -1}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0}} {
		point := Point{X: current.Player.Position.X + direction.X, Y: current.Player.Position.Y + direction.Y}
		if current.Maze.TileAt(point) == Prisoner {
			return true
		}
	}
	return false
}

// advanceGuidance expires the displayed route after 15 gameplay turns.
func (current *Game) advanceGuidance() {
	if current.GuidanceTicks > 0 {
		current.GuidanceTicks--
		if current.GuidanceTicks == 0 {
			current.GuidancePath = nil
		}
	}
}
