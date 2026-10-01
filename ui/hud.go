package ui

// This file renders the player's energy, inventory, and status feedback.

import (
	"fmt"

	"termzogs/game"

	"github.com/gdamore/tcell/v2"
)

// DrawHUD renders the energy meter, inventory, treasure status, and feedback.
func DrawHUD(screen tcell.Screen, current *game.Game) {
	width, height := screen.Size()
	if width <= 0 || height <= 0 || current == nil {
		return
	}
	if height > 1 && current.StatusMessage != "" {
		drawText(screen, 0, height-2, width, current.StatusMessage, blackBackgroundStyle().Foreground(tcell.ColorWhite))
	}
	if current.State != game.StatePlaying {
		return
	}

	capacity := current.Player.MaxEnergy
	filled := 0
	if capacity > 0 {
		filled = current.Player.Energy * 10 / capacity
	}
	bar := "██████████"
	if filled < 10 {
		bar = "██████████"[:filled*3] + "░░░░░░░░░░"[:(10-filled)*3]
	}
	holding := "NONE"
	if current.Player.HasSword {
		holding = "SWORD"
	} else if current.Player.HasTreasure {
		holding = "TREASURE"
	}
	status := fmt.Sprintf("Energy: %s %d/%d | Holding: %s | Treasure: %s | Steps: %d", bar, current.Player.Energy, capacity, holding, yesNo(current.Player.HasTreasure), current.Player.Moves)
	drawText(screen, 0, height-1, width, status, blackBackgroundStyle().Foreground(tcell.ColorWhite))
}

// yesNo formats a boolean as the HUD's inventory indicator.
func yesNo(value bool) string {
	if value {
		return "YES"
	}
	return "NO"
}
