// Package game contains player state and maze interaction rules.
package game

// Player tracks the explorer's position, inventory, and moves.
type Player struct {
	Position    Point
	HasSword    bool
	HasTreasure bool
	Moves       int
}

// Game combines the maze with the current player state.
type Game struct {
	Maze   *Maze
	Player Player
}

// NewGame generates a maze and places the player at its start.
func NewGame() *Game {
	maze := GenerateMaze()
	return &Game{Maze: maze, Player: Player{Position: maze.Start}}
}

// Move advances the player into a walkable tile and collects eligible items.
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
