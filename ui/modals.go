package ui

// This file draws title, confirmation, victory, and defeat overlays.

import (
	"fmt"

	"termzogs/game"

	"github.com/gdamore/tcell/v2"
)

// DrawOverlay renders title, item-confirmation, victory, or defeat overlays.
// It returns true when the overlay replaces the normal game view.
func DrawOverlay(screen tcell.Screen, current *game.Game) bool {
	width, height := screen.Size()
	if width <= 0 || height <= 0 || current == nil {
		return true
	}

	var lines []string
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true)
	switch current.State {
	case game.StateTitle:
		lines = []string{"TERMZOGS", "", "Find the treasure and return to the entrance.", "", "Press SPACE to begin  |  Q to quit"}
	case game.StateVictory:
		style = tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true)
		lines = []string{
			"VICTORY! YOU ESCAPED WITH THE TREASURE!",
			fmt.Sprintf("Steps Taken: %d", current.Player.Moves),
			fmt.Sprintf("Zogs Slain: %d", current.ZogsSlain),
			fmt.Sprintf("Remaining Energy: %d", current.Player.Energy),
			"Press R to play again  |  Q to quit",
		}
	case game.StateDefeat:
		style = tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)
		lines = []string{
			"GAME OVER - YOU PERISHED IN THE MAZE",
			fmt.Sprintf("Steps Taken: %d  |  Zogs Slain: %d", current.Player.Moves, current.ZogsSlain),
			"Press R to restart  |  Q to quit",
		}
	case game.StatePlaying:
		if !current.PendingSwap {
			return false
		}
		style = tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true)
		lines = []string{"Drop Treasure to pick up Sword?", "Y: drop treasure and take sword  |  N: keep treasure"}
	default:
		return false
	}

	if current.State == game.StatePlaying {
		startY := height/2 - len(lines)/2
		for index, line := range lines {
			drawCentered(screen, width, startY+index, line, style)
		}
		return false
	}
	for index, line := range lines {
		drawCentered(screen, width, height/2-len(lines)/2+index, line, style)
	}
	return true
}

// drawCentered places a clipped text line at the horizontal screen midpoint.
func drawCentered(screen tcell.Screen, width, y int, text string, style tcell.Style) {
	drawText(screen, (width-len(text))/2, y, width, text, style)
}
