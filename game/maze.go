package game

// This file defines maze data, generation, and tile queries.

import (
	"math/rand"
	"time"
)

const (
	// MazeWidth and MazeHeight define the fixed maze dimensions.
	MazeWidth  = 62
	MazeHeight = 42
)

// Point identifies a cell in the maze.
type Point struct {
	X int
	Y int
}

// Tile identifies the contents of a maze cell.
type Tile uint8

const (
	// Wall and Floor represent empty maze geometry; the remaining values are entities.
	Wall Tile = iota
	Floor
	Sword
	Prisoner
	Monster
	Treasure
	Food
)

// Maze stores the cell grid and the starting and treasure locations.
type Maze struct {
	Tiles    [MazeHeight][MazeWidth]Tile
	Start    Point
	Treasure Point
}

// GenerateMaze creates a maze using a time-based random seed.
func GenerateMaze() *Maze {
	return GenerateMazeWithSeed(time.Now().UnixNano())
}

// GenerateMazeWithSeed creates a reproducible maze and places its entities.
func GenerateMazeWithSeed(seed int64) *Maze {
	maze := &Maze{Start: Point{X: 1, Y: 1}}
	rng := rand.New(rand.NewSource(seed))

	visited := [MazeHeight][MazeWidth]bool{}
	stack := []Point{maze.Start}
	visited[maze.Start.Y][maze.Start.X] = true
	maze.Tiles[maze.Start.Y][maze.Start.X] = Floor

	directions := [...]Point{{X: 0, Y: -2}, {X: 2, Y: 0}, {X: 0, Y: 2}, {X: -2, Y: 0}}
	// Carve between unvisited cells while backtracking when a cell has no options.
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		neighbors := make([]Point, 0, len(directions))
		for _, direction := range directions {
			next := Point{X: current.X + direction.X, Y: current.Y + direction.Y}
			if next.X > 0 && next.X < MazeWidth-1 && next.Y > 0 && next.Y < MazeHeight-1 && !visited[next.Y][next.X] {
				neighbors = append(neighbors, next)
			}
		}

		if len(neighbors) == 0 {
			stack = stack[:len(stack)-1]
			continue
		}

		next := neighbors[rng.Intn(len(neighbors))]
		between := Point{X: (current.X + next.X) / 2, Y: (current.Y + next.Y) / 2}
		maze.Tiles[between.Y][between.X] = Floor
		maze.Tiles[next.Y][next.X] = Floor
		visited[next.Y][next.X] = true
		stack = append(stack, next)
	}

	// The treasure sits on the central room of the maze.
	maze.Treasure = Point{X: MazeWidth / 2, Y: MazeHeight / 2}
	maze.Tiles[maze.Treasure.Y][maze.Treasure.X] = Treasure
	place := func(tile Tile, count int) {
		for placed := 0; placed < count; {
			point := Point{X: 1 + rng.Intn(MazeWidth-2), Y: 1 + rng.Intn(MazeHeight-2)}
			if point == maze.Start || point == maze.Treasure || maze.Tiles[point.Y][point.X] != Floor {
				continue
			}
			maze.Tiles[point.Y][point.X] = tile
			placed++
		}
	}
	place(Sword, 12)
	placePrisoners(rng, maze, 35)
	place(Monster, 18)
	placeFood(rng, maze)

	return maze
}

// placePrisoners hides prisoners in wall cells adjoining a single corridor.
func placePrisoners(rng *rand.Rand, maze *Maze, count int) {
	niches := make([]Point, 0)
	directions := [...]Point{{X: 0, Y: -1}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0}}
	for y := 1; y < MazeHeight-1; y++ {
		for x := 1; x < MazeWidth-1; x++ {
			if maze.Tiles[y][x] != Wall {
				continue
			}
			openNeighbors := 0
			for _, direction := range directions {
				if maze.IsWalkable(Point{X: x + direction.X, Y: y + direction.Y}) {
					openNeighbors++
				}
			}
			if openNeighbors == 1 {
				niches = append(niches, Point{X: x, Y: y})
			}
		}
	}

	if count > len(niches) {
		count = len(niches)
	}
	for _, index := range rng.Perm(len(niches))[:count] {
		point := niches[index]
		maze.Tiles[point.Y][point.X] = Prisoner
	}
}

// placeFood scatters food across available maze dead ends.
func placeFood(rng *rand.Rand, maze *Maze) {
	deadEnds := make([]Point, 0)
	for y := 1; y < MazeHeight-1; y++ {
		for x := 1; x < MazeWidth-1; x++ {
			point := Point{X: x, Y: y}
			if maze.Tiles[y][x] != Floor || point == maze.Start {
				continue
			}
			openNeighbors := 0
			for _, direction := range [...]Point{{X: 0, Y: -1}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0}} {
				if maze.IsWalkable(Point{X: x + direction.X, Y: y + direction.Y}) {
					openNeighbors++
				}
			}
			if openNeighbors == 1 {
				deadEnds = append(deadEnds, point)
			}
		}
	}

	count := 8 + rng.Intn(5)
	if count > len(deadEnds) {
		count = len(deadEnds)
	}
	for _, index := range rng.Perm(len(deadEnds))[:count] {
		point := deadEnds[index]
		maze.Tiles[point.Y][point.X] = Food
	}
}

// TileAt returns the cell contents, treating out-of-bounds coordinates as walls.
func (maze *Maze) TileAt(point Point) Tile {
	if point.X < 0 || point.X >= MazeWidth || point.Y < 0 || point.Y >= MazeHeight {
		return Wall
	}
	return maze.Tiles[point.Y][point.X]
}

// IsWalkable reports whether a point can be entered by the player.
func (maze *Maze) IsWalkable(point Point) bool {
	tile := maze.TileAt(point)
	return tile != Wall && tile != Prisoner
}
