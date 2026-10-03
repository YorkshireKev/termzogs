package game

import (
	"math/rand"
	"time"
)

const zogPatrolRadius = 5

type zogPatrol struct {
	home     Point
	position Point
	route    []Point
	nextStep int
	active   bool
}

// UpdateZogs advances each active patrol when its movement interval has elapsed.
func (current *Game) UpdateZogs(now time.Time) {
	if current.nextZogMoveAt.IsZero() || current.State != StatePlaying || current.ExpandedView || current.CombatInProgress() {
		current.nextZogMoveAt = now.Add(ZogMoveInterval)
		return
	}
	if now.Before(current.nextZogMoveAt) {
		return
	}

	for index := range current.zogPatrols {
		patrol := &current.zogPatrols[index]
		if !patrol.active || current.Maze.TileAt(patrol.position) != Monster {
			patrol.active = false
			continue
		}
		if len(patrol.route) == 0 {
			continue
		}

		next := patrol.route[patrol.nextStep]
		if next == current.Player.Position {
			current.Maze.Tiles[patrol.position.Y][patrol.position.X] = Floor
			current.Maze.Tiles[next.Y][next.X] = Monster
			patrol.position = next
			patrol.active = false
			current.beginCombat(next)
			break
		}
		if current.Maze.TileAt(next) != Floor {
			continue
		}
		current.Maze.Tiles[patrol.position.Y][patrol.position.X] = Floor
		current.Maze.Tiles[next.Y][next.X] = Monster
		patrol.position = next
		patrol.nextStep = (patrol.nextStep + 1) % len(patrol.route)
	}
	current.nextZogMoveAt = now.Add(ZogMoveInterval)
}

func createZogPatrols(maze *Maze, seed int64) []zogPatrol {
	rng := rand.New(rand.NewSource(seed ^ 0x5a6f67506174726f))
	patrols := make([]zogPatrol, 0)
	for y, row := range maze.Tiles {
		for x, tile := range row {
			if tile != Monster {
				continue
			}
			home := Point{X: x, Y: y}
			route := createZogPatrolRoute(maze, home, rng)
			patrols = append(patrols, zogPatrol{
				home:     home,
				position: home,
				route:    route,
				active:   true,
			})
		}
	}
	return patrols
}

func createZogPatrolRoute(maze *Maze, home Point, rng *rand.Rand) []Point {
	paths := map[Point][]Point{home: {home}}
	queue := []Point{home}
	directions := [...]Point{{X: 0, Y: -1}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0}}

	for len(queue) > 0 {
		point := queue[0]
		queue = queue[1:]
		path := paths[point]
		if len(path)-1 == zogPatrolRadius {
			continue
		}
		for _, direction := range directions {
			next := Point{X: point.X + direction.X, Y: point.Y + direction.Y}
			if _, seen := paths[next]; seen || (next != home && maze.TileAt(next) != Floor) {
				continue
			}
			nextPath := append(append([]Point(nil), path...), next)
			paths[next] = nextPath
			queue = append(queue, next)
		}
	}

	candidates := make([][]Point, 0)
	furthestDistance := 0
	for y := 0; y < MazeHeight; y++ {
		for x := 0; x < MazeWidth; x++ {
			point := Point{X: x, Y: y}
			path, reachable := paths[point]
			if !reachable || point == home {
				continue
			}
			distance := len(path) - 1
			if distance >= 3 {
				candidates = append(candidates, path)
			}
			if distance > furthestDistance {
				furthestDistance = distance
			}
		}
	}
	if len(candidates) == 0 {
		for y := 0; y < MazeHeight; y++ {
			for x := 0; x < MazeWidth; x++ {
				point := Point{X: x, Y: y}
				path, reachable := paths[point]
				if reachable && point != home && len(path)-1 == furthestDistance {
					candidates = append(candidates, path)
				}
			}
		}
	}
	if len(candidates) == 0 {
		return nil
	}

	path := candidates[rng.Intn(len(candidates))]
	distance := len(path) - 1
	cycleLength := 6 + rng.Intn(3)*2
	if minimum := 2 * distance; cycleLength < minimum {
		cycleLength = minimum
	}
	extraSteps := cycleLength - 2*distance
	route := make([]Point, 0, cycleLength)
	route = append(route, path[1:]...)
	for step := 0; step < extraSteps; step += 2 {
		route = append(route, path[distance-1], path[distance])
	}
	for index := distance - 1; index >= 0; index-- {
		route = append(route, path[index])
	}
	return route
}
