package ui

import (
	"testing"

	"termzogs/game"

	"github.com/gdamore/tcell/v2"
)

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
	viewLeft := (width - ViewportSize) / 2
	viewTop := (screenHeight - ViewportSize) / 2
	worldLeft := current.Player.Position.X - ViewportSize/2
	worldTop := current.Player.Position.Y - ViewportSize/2
	for point, want := range map[game.Point]rune{
		{X: 11, Y: 10}: '↓',
		{X: 11, Y: 11}: '↓',
		{X: 11, Y: 12}: '↓',
	} {
		x := viewLeft + point.X - worldLeft
		y := viewTop + point.Y - worldTop
		if got := cells[y*width+x].Runes[0]; got != want {
			t.Errorf("rendered path cell %v = %q, want %q", point, got, want)
		}
	}
}
