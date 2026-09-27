package ui

import (
	"fmt"

	"termzogs/game"

	"github.com/gdamore/tcell/v2"
)

const ViewportSize = 15

type Renderer struct{}

func (Renderer) Draw(screen tcell.Screen, current *game.Game) {
	screenWidth, screenHeight := screen.Size()
	screen.Clear()
	if screenWidth <= 0 || screenHeight <= 0 {
		return
	}

	viewWidth := ViewportSize
	if viewWidth > screenWidth {
		viewWidth = screenWidth
	}
	viewHeight := ViewportSize
	if viewHeight > screenHeight {
		viewHeight = screenHeight
	}
	viewLeft := (screenWidth - viewWidth) / 2
	viewTop := (screenHeight - viewHeight) / 2
	worldLeft := current.Player.Position.X - viewWidth/2
	worldTop := current.Player.Position.Y - viewHeight/2

	for y := 0; y < viewHeight; y++ {
		for x := 0; x < viewWidth; x++ {
			world := game.Point{X: worldLeft + x, Y: worldTop + y}
			tile := current.Maze.TileAt(world)
			character, style := tileAppearance(tile)
			if world == current.Player.Position {
				character = '@'
				style = tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true)
			}
			screen.SetContent(viewLeft+x, viewTop+y, character, nil, style)
		}
	}

	if screenHeight > viewTop+viewHeight {
		status := fmt.Sprintf("Moves: %d  Sword: %t  Treasure: %t  (Q to quit)", current.Player.Moves, current.Player.HasSword, current.Player.HasTreasure)
		drawText(screen, (screenWidth-len(status))/2, viewTop+viewHeight, screenWidth, status, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	}
}

func tileAppearance(tile game.Tile) (rune, tcell.Style) {
	switch tile {
	case game.Floor:
		return ' ', tcell.StyleDefault.Foreground(tcell.ColorBlack)
	case game.Sword:
		return '/', tcell.StyleDefault.Foreground(tcell.ColorAqua)
	case game.Prisoner:
		return 'P', tcell.StyleDefault.Foreground(tcell.ColorGreen)
	case game.Monster:
		return 'M', tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)
	case game.Treasure:
		return '$', tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true)
	default:
		return '█', tcell.StyleDefault.Foreground(tcell.ColorGray)
	}
}

func drawText(screen tcell.Screen, x, y, screenWidth int, text string, style tcell.Style) {
	for _, character := range text {
		if x >= screenWidth {
			return
		}
		if x >= 0 {
			screen.SetContent(x, y, character, nil, style)
		}
		x++
	}
}
