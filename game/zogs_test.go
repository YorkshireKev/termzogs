package game

import (
	"testing"
	"time"
)

func TestZogsMoveOneCellPerSecondInLocalPatrolLoops(t *testing.T) {
	current := NewGameWithSeed(31)
	if len(current.zogPatrols) == 0 {
		t.Fatal("game did not create any Zog patrols")
	}

	for _, patrol := range current.zogPatrols {
		if len(patrol.route) < 6 || len(patrol.route) > 10 {
			t.Fatalf("patrol at %v has %d steps, want 6 through 10", patrol.home, len(patrol.route))
		}
		if patrol.route[len(patrol.route)-1] != patrol.home {
			t.Fatalf("patrol at %v ends at %v instead of returning home", patrol.home, patrol.route[len(patrol.route)-1])
		}
		for _, point := range patrol.route {
			distance := absInt(point.X-patrol.home.X) + absInt(point.Y-patrol.home.Y)
			if distance > zogPatrolRadius || (point != patrol.home && current.Maze.TileAt(point) != Floor) {
				t.Fatalf("patrol at %v leaves its starting area or enters a non-floor cell at %v", patrol.home, point)
			}
		}
	}

	patrol := &current.zogPatrols[0]
	start := current.nextZogMoveAt
	initialPosition := patrol.position
	current.UpdateZogs(start.Add(-time.Nanosecond))
	if patrol.position != initialPosition {
		t.Fatal("Zog moved before one second elapsed")
	}
	current.UpdateZogs(start)
	if patrol.position != patrol.route[0] {
		t.Fatalf("Zog moved to %v, want first patrol cell %v", patrol.position, patrol.route[0])
	}
	current.UpdateZogs(start.Add(ZogMoveInterval - time.Nanosecond))
	if patrol.position != patrol.route[0] {
		t.Fatal("Zog moved more than once during a one-second interval")
	}
}

func TestExpandedViewPausesZogPatrols(t *testing.T) {
	current := NewGameWithSeed(32)
	patrol := &current.zogPatrols[0]
	initialPosition := patrol.position
	pauseAt := current.nextZogMoveAt
	current.ExpandedView = true
	current.UpdateZogs(pauseAt)
	if patrol.position != initialPosition {
		t.Fatal("Zog moved while the expanded map was open")
	}
	resumedAt := pauseAt.Add(10 * time.Second)
	current.UpdateZogs(resumedAt)
	if patrol.position != initialPosition {
		t.Fatal("Zog moved during an extended expanded-map display")
	}

	current.ExpandedView = false
	current.UpdateZogs(resumedAt.Add(ZogMoveInterval - time.Nanosecond))
	if patrol.position != initialPosition {
		t.Fatal("Zog moved before a full second after the expanded map closed")
	}
	current.UpdateZogs(resumedAt.Add(ZogMoveInterval))
	if patrol.position != patrol.route[0] {
		t.Fatalf("Zog moved to %v after resuming, want %v", patrol.position, patrol.route[0])
	}
}

func TestZogMovingOntoPlayerStartsCombat(t *testing.T) {
	current := NewGameWithSeed(33)
	patrol := &current.zogPatrols[0]
	playerPosition := current.Player.Position
	current.Maze.Tiles[patrol.position.Y][patrol.position.X] = Floor
	current.Maze.Tiles[playerPosition.Y][playerPosition.X] = Floor
	patrol.position = patrol.home
	current.Maze.Tiles[patrol.home.Y][patrol.home.X] = Monster
	patrol.route = []Point{playerPosition, patrol.home}
	patrol.nextStep = 0
	patrol.active = true
	current.Player.HasSword = true
	moveAt := current.nextZogMoveAt

	current.UpdateZogs(moveAt)

	if !current.CombatInProgress() || current.combatPoint != playerPosition {
		t.Fatalf("combat=%t point=%v; want combat at player position %v", current.CombatInProgress(), current.combatPoint, playerPosition)
	}
	if current.Maze.TileAt(playerPosition) != Monster || current.Maze.TileAt(patrol.home) != Floor {
		t.Fatal("Zog did not move from its home onto the player")
	}

	if !current.ResolveCombatIfReady(current.combatStartedAt.Add(CombatAnimationDuration)) {
		t.Fatal("combat did not resolve")
	}
	if current.ZogsSlain != 1 || current.Maze.TileAt(playerPosition) != Floor || current.Player.HasSword {
		t.Fatal("moving Zog encounter did not resolve using normal armed-combat rules")
	}
}
