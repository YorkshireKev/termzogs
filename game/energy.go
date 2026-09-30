package game

// This file applies movement, combat, and food energy changes.

// spendEnergy applies a cost and immediately defeats a player at zero energy.
func (game *Game) spendEnergy(amount int) {
	if game.State != StatePlaying {
		return
	}
	game.Player.Energy -= amount
	if game.Player.Energy <= 0 {
		game.Player.Energy = 0
		game.State = StateDefeat
		game.SetStatusMessage("You ran out of energy.")
	}
}

// RestoreEnergy adds energy without exceeding the player's capacity.
func (game *Game) RestoreEnergy(amount int) {
	game.Player.Energy += amount
	if game.Player.Energy > game.Player.MaxEnergy {
		game.Player.Energy = game.Player.MaxEnergy
	}
}
