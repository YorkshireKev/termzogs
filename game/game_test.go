package game

import "testing"

func TestNewGameStartAndRestartMazeLifecycle(t *testing.T) {
	current := NewTitleGame()
	initialSeed := current.seed

	current.StartNewGame()
	if current.seed == initialSeed {
		t.Fatal("starting a new game reused the previous seed")
	}
	firstMaze := *current.Maze
	firstSeed := current.seed

	current.State = StateDefeat
	current.Restart()
	if current.seed != firstSeed {
		t.Fatalf("restart changed seed from %d to %d", firstSeed, current.seed)
	}
	if *current.Maze != firstMaze {
		t.Fatal("restart generated a different maze")
	}

	current.StartNewGame()
	if current.seed == firstSeed {
		t.Fatal("starting another game reused the previous maze seed")
	}
	if *current.Maze == firstMaze {
		t.Fatal("starting another game reused the previous maze")
	}
}
