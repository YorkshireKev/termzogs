package game

// This file finds shortest maze routes for prisoner guidance.

// FindPath returns the shortest walkable route from start to goal, excluding start.
func (maze *Maze) FindPath(start, goal Point) ([]Point, bool) {
	if !maze.IsWalkable(start) || !maze.IsWalkable(goal) {
		return nil, false
	}
	if start == goal {
		return []Point{}, true
	}

	queue := []Point{start}
	previous := map[Point]Point{}
	visited := map[Point]bool{start: true}
	directions := [...]Point{{X: 0, Y: -1}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0}}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, direction := range directions {
			next := Point{X: current.X + direction.X, Y: current.Y + direction.Y}
			if visited[next] || !maze.IsWalkable(next) {
				continue
			}
			visited[next] = true
			previous[next] = current
			if next == goal {
				path := []Point{goal}
				for path[len(path)-1] != start {
					path = append(path, previous[path[len(path)-1]])
				}
				for left, right := 0, len(path)-1; left < right; left, right = left+1, right-1 {
					path[left], path[right] = path[right], path[left]
				}
				return path[1:], true
			}
			queue = append(queue, next)
		}
	}
	return nil, false
}
