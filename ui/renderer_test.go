package ui

import (
	"testing"

	"termzogs/game"

	"github.com/gdamore/tcell/v2"
)

func TestLocalCellPatternsMatchViewportCellSize(t *testing.T) {
	for character, pattern := range localCellPatterns {
		for rowIndex, row := range pattern {
			if width := len([]rune(row)); width != ViewportCellWidth {
				t.Errorf("pattern %q row %d has width %d, want %d", character, rowIndex, width, ViewportCellWidth)
			}
		}
	}
}

func TestDrawAlwaysUsesBlackBackground(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	if err := screen.Init(); err != nil {
		t.Fatal(err)
	}
	defer screen.Fini()

	current := &game.Game{State: game.StateTitle}
	Renderer{}.Draw(screen, current)
	screen.Show()
	cells, _, _ := screen.GetContents()
	for index, cell := range cells {
		_, background, _ := cell.Style.Decompose()
		if background != tcell.ColorBlack {
			t.Errorf("cell %d background = %v, want black", index, background)
		}
	}
}

func TestDrawUsesCustomPatternInLocalView(t *testing.T) {
	originalPattern := localCellPatterns['M']
	localCellPatterns['M'] = cellPattern{"ABCDEF", "GHIJKL", "MNOPQR"}
	defer func() {
		localCellPatterns['M'] = originalPattern
	}()

	screen := tcell.NewSimulationScreen("")
	if err := screen.Init(); err != nil {
		t.Fatal(err)
	}
	defer screen.Fini()

	current := &game.Game{
		Maze: &game.Maze{},
		Player: game.Player{
			Position:  game.Point{X: 10, Y: 10},
			MaxEnergy: game.MaxEnergy,
			Energy:    game.InitialEnergy,
		},
		State: game.StatePlaying,
	}
	monster := game.Point{X: 11, Y: 10}
	current.Maze.Tiles[monster.Y][monster.X] = game.Monster

	Renderer{}.Draw(screen, current)
	screen.Show()
	cells, width, height := screen.GetContents()
	viewLeft := (width - ViewportWidth*ViewportCellWidth) / 2
	viewTop := (height - 2 - ViewportHeight*ViewportCellHeight) / 2
	worldLeft := current.Player.Position.X - ViewportWidth/2
	worldTop := current.Player.Position.Y - ViewportHeight/2
	cellLeft := viewLeft + (monster.X-worldLeft)*ViewportCellWidth
	cellTop := viewTop + (monster.Y-worldTop)*ViewportCellHeight
	for row, expected := range localCellPatterns['M'] {
		for column, want := range []rune(expected) {
			if got := cells[(cellTop+row)*width+cellLeft+column].Runes[0]; got != want {
				t.Errorf("rendered pattern at (%d, %d) = %q, want %q", column, row, got, want)
			}
		}
	}
}

func TestDrawRendersGuidanceArrowsAlongPathTurns(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	if err := screen.Init(); err != nil {
		t.Fatal(err)
	}
	defer screen.Fini()

	current := &game.Game{
		Maze: &game.Maze{},
		Player: game.Player{
			Position:  game.Point{X: 10, Y: 10},
			MaxEnergy: game.MaxEnergy,
			Energy:    game.InitialEnergy,
		},
		State:         game.StatePlaying,
		GuidanceTicks: 15,
		GuidancePath: []game.Point{
			{X: 11, Y: 10},
			{X: 11, Y: 11},
			{X: 11, Y: 12},
		},
	}

	Renderer{}.Draw(screen, current)
	screen.Show()
	cells, width, _ := screen.GetContents()
	_, screenHeight := screen.Size()
	viewLeft := (width - ViewportWidth*ViewportCellWidth) / 2
	viewTop := (screenHeight - 2 - ViewportHeight*ViewportCellHeight) / 2
	worldLeft := current.Player.Position.X - ViewportWidth/2
	worldTop := current.Player.Position.Y - ViewportHeight/2
	for point, want := range map[game.Point]rune{
		{X: 11, Y: 10}: '↓',
		{X: 11, Y: 11}: '↓',
		{X: 11, Y: 12}: '↓',
	} {
		x := viewLeft + (point.X-worldLeft)*ViewportCellWidth
		y := viewTop + (point.Y-worldTop)*ViewportCellHeight
		if got := cells[y*width+x].Runes[0]; got != want {
			t.Errorf("rendered path cell %v = %q, want %q", point, got, want)
		}
		for row := 0; row < ViewportCellHeight; row++ {
			for column := 0; column < ViewportCellWidth; column++ {
				if got := cells[(y+row)*width+x+column].Runes[0]; got != want {
					t.Errorf("rendered path block %v at (%d, %d) = %q, want %q", point, column, row, got, want)
				}
			}
		}
	}
}

func TestDrawUsesExpandedMazeViewport(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	if err := screen.Init(); err != nil {
		t.Fatal(err)
	}
	defer screen.Fini()

	current := &game.Game{
		Maze: &game.Maze{},
		Player: game.Player{
			Position:  game.Point{X: 31, Y: 20},
			MaxEnergy: game.MaxEnergy,
			Energy:    game.InitialEnergy,
		},
		State:         game.StatePlaying,
		ExpandedView:  true,
		StatusMessage: "Status text",
	}
	wideEdge := game.Point{X: 61, Y: 20}
	current.Maze.Tiles[wideEdge.Y][wideEdge.X] = game.Food

	Renderer{}.Draw(screen, current)
	screen.Show()
	cells, screenWidth, screenHeight := screen.GetContents()
	viewLeft := (screenWidth - ExpandedViewportWidth) / 2
	viewHeight := ExpandedViewportHeight
	if viewHeight > screenHeight-3 {
		viewHeight = screenHeight - 3
	}
	viewTop := (screenHeight-viewHeight)/2 - 1
	if viewTop < 0 {
		viewTop = 0
	}
	worldLeft := current.Player.Position.X - ExpandedViewportWidth/2
	worldTop := current.Player.Position.Y - viewHeight/2
	cellX := viewLeft + wideEdge.X - worldLeft
	cellY := viewTop + wideEdge.Y - worldTop
	if got := cells[cellY*screenWidth+cellX].Runes[0]; got != '%' {
		t.Fatalf("wide edge rendered %q, want food marker", got)
	}
	mapBottom := viewTop + viewHeight
	if got := cells[mapBottom*screenWidth+screenWidth/2].Runes[0]; got != ' ' {
		t.Fatalf("map separator row = %q, want blank", got)
	}
	if got := cells[(screenHeight-2)*screenWidth].Runes[0]; got != 'S' {
		t.Fatalf("status message rendered %q, want 'S'", got)
	}
}
