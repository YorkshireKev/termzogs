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

func TestTitleInstructionsOpenAndAnyKeyReturns(t *testing.T) {
	for _, key := range []rune{'i', 'I'} {
		current := game.NewTitleGame()
		handleKey(current, tcell.NewEventKey(tcell.KeyRune, key, 0))
		if current.State != game.StateInstructions {
			t.Fatalf("key %q set state to %v, want instructions", key, current.State)
		}

		for _, event := range []*tcell.EventKey{
			tcell.NewEventKey(tcell.KeyRune, 'q', 0),
			tcell.NewEventKey(tcell.KeyEscape, 0, 0),
			tcell.NewEventKey(tcell.KeyRune, ' ', 0),
		} {
			if quit := handleKey(current, event); quit {
				t.Fatalf("key %v quit instead of returning to title", event)
			}
			if current.State != game.StateTitle {
				t.Fatalf("key %v set state to %v, want title", event, current.State)
			}
			current.State = game.StateInstructions
		}
	}
}

func TestGameOverReturnsToTitleAndCanStartAgain(t *testing.T) {
	for _, key := range []rune{'m', 'M'} {
		current := game.NewGameWithSeed(23)
		current.State = game.StateDefeat

		if quit := handleKey(current, tcell.NewEventKey(tcell.KeyRune, key, 0)); quit {
			t.Fatalf("key %q unexpectedly quit the game", key)
		}
		if current.State != game.StateTitle {
			t.Fatalf("key %q set state to %v, want title", key, current.State)
		}

		handleKey(current, tcell.NewEventKey(tcell.KeyRune, ' ', 0))
		if current.State != game.StatePlaying {
			t.Fatalf("space after returning to title set state to %v, want playing", current.State)
		}
		if current.Maze == nil {
			t.Fatal("space after returning to title did not generate a maze")
		}
	}
}

func TestPlayingQuitKeysReturnToTitle(t *testing.T) {
	for _, event := range []*tcell.EventKey{
		tcell.NewEventKey(tcell.KeyRune, 'q', 0),
		tcell.NewEventKey(tcell.KeyRune, 'Q', 0),
		tcell.NewEventKey(tcell.KeyEscape, 0, 0),
	} {
		current := game.NewGameWithSeed(27)
		if quit := handleKey(current, event); quit {
			t.Fatalf("key %v exited the program instead of returning to title", event)
		}
		if current.State != game.StateTitle {
			t.Fatalf("key %v set state to %v, want title", event, current.State)
		}
	}
}

func TestQuitKeysExitOnlyFromTitleScreen(t *testing.T) {
	for _, event := range []*tcell.EventKey{
		tcell.NewEventKey(tcell.KeyRune, 'q', 0),
		tcell.NewEventKey(tcell.KeyRune, 'Q', 0),
		tcell.NewEventKey(tcell.KeyEscape, 0, 0),
	} {
		current := game.NewTitleGame()
		if !handleKey(current, event) {
			t.Errorf("key %v on the title screen did not exit the program", event)
		}

		for _, state := range []game.GameState{
			game.StateInstructions,
			game.StatePlaying,
			game.StateDefeat,
			game.StateVictory,
		} {
			current := game.NewTitleGame()
			current.State = state
			if quit := handleKey(current, event); quit {
				t.Errorf("key %v in state %v exited the program", event, state)
			}
			if current.State != game.StateTitle {
				t.Errorf("key %v in state %v set state to %v, want title", event, state, current.State)
			}
		}
	}
}

func TestWebQuitKeysDoNotExitFromTitleScreen(t *testing.T) {
	for _, event := range []*tcell.EventKey{
		tcell.NewEventKey(tcell.KeyRune, 'q', 0),
		tcell.NewEventKey(tcell.KeyRune, 'Q', 0),
		tcell.NewEventKey(tcell.KeyEscape, 0, 0),
		tcell.NewEventKey(tcell.KeyCtrlC, 0, 0),
	} {
		current := game.NewTitleGame()
		if handleKeyForPlatform(current, event, true) {
			t.Errorf("key %v exited the web game from the title screen", event)
		}
		if current.State != game.StateTitle {
			t.Errorf("key %v set state to %v, want title", event, current.State)
		}
	}
}

func TestWebQuitKeysReturnToTitleFromGameplay(t *testing.T) {
	for _, event := range []*tcell.EventKey{
		tcell.NewEventKey(tcell.KeyRune, 'q', 0),
		tcell.NewEventKey(tcell.KeyRune, 'Q', 0),
		tcell.NewEventKey(tcell.KeyEscape, 0, 0),
	} {
		current := game.NewGameWithSeed(28)
		if handleKeyForPlatform(current, event, true) {
			t.Errorf("key %v exited the web game instead of returning to title", event)
		}
		if current.State != game.StateTitle {
			t.Errorf("key %v set state to %v, want title", event, current.State)
		}
	}
}

func TestCombatBlocksGameplayInputButAllowsReturningToTitle(t *testing.T) {
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
	if handleKey(current, tcell.NewEventKey(tcell.KeyRune, 'q', 0)) || current.State != game.StateTitle {
		t.Fatal("quit input during combat did not return to title")
	}
}
