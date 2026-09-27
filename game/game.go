package game

type Player struct {
	Position    Point
	HasSword    bool
	HasTreasure bool
	Moves       int
}

type Game struct {
	Maze   *Maze
	Player Player
}

func NewGame() *Game {
	maze := GenerateMaze()
	return &Game{Maze: maze, Player: Player{Position: maze.Start}}
}

func (game *Game) Move(dx, dy int) {
	next := Point{X: game.Player.Position.X + dx, Y: game.Player.Position.Y + dy}
	if !game.Maze.IsWalkable(next) {
		return
	}

	game.Player.Position = next
	game.Player.Moves++
	switch game.Maze.TileAt(next) {
	case Sword:
		if !game.Player.HasTreasure {
			game.Player.HasSword = true
			game.Maze.Tiles[next.Y][next.X] = Floor
		}
	case Treasure:
		if !game.Player.HasSword {
			game.Player.HasTreasure = true
			game.Maze.Tiles[next.Y][next.X] = Floor
		}
	}
}
