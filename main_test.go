package main

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"termzogs/game"
)

func TestNextMovementKeyDismissesExpandedViewAndMoves(t *testing.T) {
	current := game.NewGameWithSeed(21)
	current.Player.Position = game.Point{X: 2, Y: 2}
	current.Maze.Tiles[2][2] = game.Floor
	current.Maze.Tiles[2][1] = game.Floor
	current.ExpandedView = true

	if quit := handleKey(current, tcell.NewEventKey(tcell.KeyRune, 'a', 0)); quit {
		t.Fatal("movement key unexpectedly quit the game")
	}
	if current.ExpandedView || current.Player.Position != (game.Point{X: 1, Y: 2}) {
		t.Fatalf("expanded=%t position=%v, want local view and movement to (1, 2)", current.ExpandedView, current.Player.Position)
	}
}

func TestVKeyDismissesExpandedViewWithoutReopeningOrCharging(t *testing.T) {
	current := game.NewGameWithSeed(22)
	current.Player.Energy = 80
	current.ExpandedView = true

	handleKey(current, tcell.NewEventKey(tcell.KeyRune, 'v', 0))
	if current.ExpandedView || current.Player.Energy != 80 {
		t.Fatalf("expanded=%t energy=%d, want closed view without another charge", current.ExpandedView, current.Player.Energy)
	}
}
