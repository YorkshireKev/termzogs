package game

// This file resolves sword-based and unarmed encounters with Zogs.

import (
	"math/rand"
	"time"
)

func (game *Game) CombatInProgress() bool {
	return !game.combatStartedAt.IsZero()
}

func (game *Game) beginCombat(point Point) {
	game.combatPoint = point
	game.combatStartedAt = time.Now()
	game.SetStatusMessage("Fighting the Zog...")
}

// ResolveCombatIfReady reveals the fight outcome after its animation duration.
func (game *Game) ResolveCombatIfReady(now time.Time) bool {
	if !game.CombatInProgress() || now.Sub(game.combatStartedAt) < CombatAnimationDuration {
		return false
	}

	point := game.combatPoint
	game.combatStartedAt = time.Time{}
	game.combatPoint = Point{}
	game.resolveCombat(point)
	if game.State == StatePlaying && game.Player.Position == game.Maze.Start && game.Player.HasTreasure {
		game.State = StateVictory
		game.SetStatusMessage("You escaped with the treasure!")
	}
	return true
}

// resolveCombat applies armed or 50/50 unarmed combat against a Zog.
func (game *Game) resolveCombat(point Point) {
	if game.Player.HasSword {
		game.Player.HasSword = false
		game.Maze.Tiles[point.Y][point.X] = Floor
		game.ZogsSlain++
		game.spendEnergy(5)
		game.SetStatusMessage("You slew the Zog, but your sword shattered!")
		return
	}

	if game.rng == nil {
		game.rng = rand.New(rand.NewSource(1))
	}
	if game.rng.Float32() < 0.5 {
		game.Maze.Tiles[point.Y][point.X] = Floor
		game.ZogsSlain++
		game.spendEnergy(40)
		game.SetStatusMessage("Desperate fight! You defeated the Zog but lost 40 Energy!")
		return
	}
	game.State = StateDefeat
	game.DefeatCause = DefeatCauseZog
	game.SetStatusMessage("The Zog overpowered you!")
}
