package game

// This file resolves sword-based and unarmed encounters with Zogs.

import "math/rand"

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
	game.SetStatusMessage("The Zog overpowered you!")
}
