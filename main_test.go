package main

import (
	"testing"

	"termzogs/game"

	"github.com/gdamore/tcell/v2"
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

func TestMapKeysOpenAndCloseExpandedView(t *testing.T) {
	for _, key := range []rune{'v', 'V', 'm', 'M'} {
		current := game.NewGameWithSeed(22)
		initialEnergy := current.Player.Energy

		handleKey(current, tcell.NewEventKey(tcell.KeyRune, key, 0))
		if !current.ExpandedView || current.Player.Energy != initialEnergy-game.MazeViewCost {
			t.Fatalf("key %q: expanded=%t energy=%d, want map open and one charge", key, current.ExpandedView, current.Player.Energy)
		}

		handleKey(current, tcell.NewEventKey(tcell.KeyRune, key, 0))
		if current.ExpandedView || current.Player.Energy != initialEnergy-game.MazeViewCost {
			t.Fatalf("key %q: expanded=%t energy=%d, want map closed without another charge", key, current.ExpandedView, current.Player.Energy)
		}
	}
}

func TestCombatBlocksGameplayInputButAllowsQuit(t *testing.T) {
	current := game.NewGameWithSeed(26)
	current.Maze.Tiles[1][2] = game.Monster
	handleKey(current, tcell.NewEventKey(tcell.KeyRight, 0, 0))
	if !current.CombatInProgress() {
		t.Fatal("stepping onto a Zog did not start combat")
	}
	position := current.Player.Position
	moves := current.Player.Moves

	handleKey(current, tcell.NewEventKey(tcell.KeyDown, 0, 0))
	if current.Player.Position != position || current.Player.Moves != moves {
		t.Fatal("movement input was accepted during combat")
	}
	if !handleKey(current, tcell.NewEventKey(tcell.KeyRune, 'q', 0)) {
		t.Fatal("quit input was blocked during combat")
	}
}
