package game

// This file verifies maze boundaries, connectivity, and generated entity counts.

import "testing"

// TestGeneratedMazeHasClosedBoundaryAndConnectedPassages checks maze invariants.
func TestGeneratedMazeHasClosedBoundaryAndConnectedPassages(t *testing.T) {
	maze := GenerateMazeWithSeed(7)
	for x := 0; x < MazeWidth; x++ {
		if maze.Tiles[0][x] != Wall || maze.Tiles[MazeHeight-1][x] != Wall {
			t.Fatalf("horizontal boundary opened at x=%d", x)
		}
	}
	for y := 0; y < MazeHeight; y++ {
		if maze.Tiles[y][0] != Wall || maze.Tiles[y][MazeWidth-1] != Wall {
			t.Fatalf("vertical boundary opened at y=%d", y)
		}
	}

	seen := map[Point]bool{maze.Start: true}
	queue := []Point{maze.Start}
	directions := [...]Point{{X: 0, Y: -1}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0}}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, direction := range directions {
			next := Point{X: current.X + direction.X, Y: current.Y + direction.Y}
			if !seen[next] && maze.IsWalkable(next) {
				seen[next] = true
				queue = append(queue, next)
			}
		}
	}

	counts := map[Tile]int{}
	for y, row := range maze.Tiles {
		for x, tile := range row {
			if tile != Wall && !seen[Point{X: x, Y: y}] {
				t.Fatalf("walkable tile at (%d, %d) is unreachable", x, y)
			}
			counts[tile]++
			if tile == Food {
				openNeighbors := 0
				for _, direction := range directions {
					neighbor := Point{X: x + direction.X, Y: y + direction.Y}
					if maze.IsWalkable(neighbor) {
						openNeighbors++
					}
				}
				if openNeighbors != 1 {
					t.Errorf("food at (%d, %d) has %d open neighbors, want dead end", x, y, openNeighbors)
				}
			}
		}
	}
	for tile, want := range map[Tile]int{Sword: 12, Prisoner: 5, Monster: 18, Treasure: 1} {
		if counts[tile] != want {
			t.Errorf("tile %d count = %d, want %d", tile, counts[tile], want)
		}
	}
	if counts[Food] < 8 || counts[Food] > 12 {
		t.Errorf("food count = %d, want 8 through 12", counts[Food])
	}
}
