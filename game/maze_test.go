package game

// This file verifies maze boundaries, connectivity, and generated entity counts.

import "testing"

// TestGeneratedMazeHasClosedBoundaryAndConnectedPassages checks maze invariants.
func TestGeneratedMazeHasClosedBoundaryAndConnectedPassages(t *testing.T) {
	maze := GenerateMazeWithSeed(7)
	route, found := maze.FindPath(maze.Start, maze.Treasure)
	if !found {
		t.Fatal("treasure is unreachable")
	}
	mainRoute := make(map[Point]bool, len(route))
	for _, point := range route {
		mainRoute[point] = true
	}
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
			if maze.IsWalkable(Point{X: x, Y: y}) && !seen[Point{X: x, Y: y}] {
				t.Fatalf("walkable tile at (%d, %d) is unreachable", x, y)
			}
			counts[tile]++
			if tile == Prisoner {
				point := Point{X: x, Y: y}
				openNeighbors := 0
				for _, direction := range directions {
					if maze.IsWalkable(Point{X: x + direction.X, Y: y + direction.Y}) {
						openNeighbors++
					}
				}
				if maze.IsWalkable(point) || openNeighbors != 1 {
					t.Errorf("prisoner at (%d, %d) is walkable or not in a wall niche", x, y)
				}
			}
			if tile == Food {
				openNeighbors := 0
				for _, direction := range directions {
					neighbor := Point{X: x + direction.X, Y: y + direction.Y}
					if maze.IsWalkable(neighbor) {
						openNeighbors++
					}
				}
				if openNeighbors != 1 && !mainRoute[Point{X: x, Y: y}] {
					t.Errorf("food at (%d, %d) has %d open neighbors and is not on the safe route", x, y, openNeighbors)
				}
			}
		}
	}
	for tile, want := range map[Tile]int{Sword: 18, Prisoner: 35, Monster: 18, Treasure: 1} {
		if counts[tile] != want {
			t.Errorf("tile %d count = %d, want %d", tile, counts[tile], want)
		}
	}
	if counts[Food] < 15 || counts[Food] > 20 {
		t.Errorf("food count = %d, want 15 through 20", counts[Food])
	}
}

func TestGeneratedMazeSwordAndFoodCountsAcrossSeeds(t *testing.T) {
	if InitialEnergy != 250 || FoodPlacementBaselineEnergy != 200 {
		t.Fatalf("starting energy=%d and food-placement baseline=%d; want 250 and 200", InitialEnergy, FoodPlacementBaselineEnergy)
	}

	for seed := int64(0); seed < 100; seed++ {
		maze := GenerateMazeWithSeed(seed)
		counts := map[Tile]int{}
		for _, row := range maze.Tiles {
			for _, tile := range row {
				counts[tile]++
			}
		}
		if counts[Sword] != 18 {
			t.Errorf("seed %d has %d swords, want 18", seed, counts[Sword])
		}
		if counts[Food] < 15 || counts[Food] > 20 {
			t.Errorf("seed %d has %d food pickups, want 15 through 20", seed, counts[Food])
		}
	}
}

func TestGeneratedMazeHasAlternateRoutesAndOpenAreas(t *testing.T) {
	for seed := int64(0); seed < 20; seed++ {
		maze := GenerateMazeWithSeed(seed)
		path, found := maze.FindPath(maze.Start, maze.Treasure)
		if !found {
			t.Fatalf("seed %d has no route to treasure", seed)
		}

		hasAlternateRoute := false
		for _, point := range path[1 : len(path)-1] {
			original := maze.Tiles[point.Y][point.X]
			maze.Tiles[point.Y][point.X] = Wall
			_, alternateFound := maze.FindPath(maze.Start, maze.Treasure)
			maze.Tiles[point.Y][point.X] = original
			if alternateFound {
				hasAlternateRoute = true
				break
			}
		}
		if !hasAlternateRoute {
			t.Fatalf("seed %d has no alternate route to treasure", seed)
		}

		hasOpenArea := false
		for y := 1; y < MazeHeight-3 && !hasOpenArea; y++ {
			for x := 1; x < MazeWidth-3 && !hasOpenArea; x++ {
				allWalkable := true
				for roomY := y; roomY < y+3; roomY++ {
					for roomX := x; roomX < x+3; roomX++ {
						if !maze.IsWalkable(Point{X: roomX, Y: roomY}) {
							allWalkable = false
						}
					}
				}
				hasOpenArea = allWalkable
			}
		}
		if !hasOpenArea {
			t.Fatalf("seed %d has no open 3x3 area", seed)
		}
	}
}

func TestGeneratedSafeRouteCanBeCompletedWithEnergy(t *testing.T) {
	for seed := int64(0); seed < 100; seed++ {
		maze := GenerateMazeWithSeed(seed)
		path, found := maze.FindPath(maze.Start, maze.Treasure)
		if !found {
			t.Fatalf("seed %d has no route to treasure", seed)
		}

		capacity := MaxEnergy
		if len(path)+2 > capacity {
			capacity = len(path) + 2
		}
		current := NewGameWithSeed(seed)
		if current.Player.Energy != 250 {
			t.Fatalf("seed %d starts with %d energy, want 250", seed, current.Player.Energy)
		}
		if current.Player.MaxEnergy != capacity {
			t.Fatalf("seed %d capacity = %d, want %d", seed, current.Player.MaxEnergy, capacity)
		}
		energy := InitialEnergy
		hasSword := false
		routeZogs := 0
		trip := append([]Point(nil), path...)
		for index := len(path) - 2; index >= 0; index-- {
			trip = append(trip, path[index])
		}
		trip = append(trip, maze.Start)
		for step, point := range trip {
			tile := maze.TileAt(point)
			if tile == Sword {
				hasSword = true
			}
			energy--
			if tile == Monster {
				routeZogs++
				if !hasSword {
					t.Fatalf("seed %d reaches a route Zog without a sword", seed)
				}
				energy -= 5
				hasSword = false
				maze.Tiles[point.Y][point.X] = Floor
			}
			if energy <= 0 {
				t.Fatalf("seed %d ran out of energy at step %d/%d at (%d, %d), route length %d, capacity %d", seed, step+1, len(trip), point.X, point.Y, len(path), capacity)
			}
			if maze.TileAt(point) == Food {
				energy += 50
				if energy > capacity {
					energy = capacity
				}
				maze.Tiles[point.Y][point.X] = Floor
			}
		}
		if routeZogs != 2 {
			t.Fatalf("seed %d has %d Zogs on the main route, want 2", seed, routeZogs)
		}
	}
}
