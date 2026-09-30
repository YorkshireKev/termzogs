package ui

// This file draws the player-centered maze viewport and status line.

import (
	"termzogs/game"

	"github.com/gdamore/tcell/v2"
)

// ViewportSize is the width and height of the local gameplay view.
const ViewportSize = 15

// Renderer draws game state onto a tcell screen.
type Renderer struct{}

// Draw paints the local viewport and player status.
func (Renderer) Draw(screen tcell.Screen, current *game.Game) {
	screenWidth, screenHeight := screen.Size()
	screen.Clear()
	if screenWidth <= 0 || screenHeight <= 0 {
		return
	}
	if DrawOverlay(screen, current) {
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
	guidance := make(map[game.Point]rune)
	if current.GuidanceTicks > 0 {
		for index, point := range current.GuidancePath {
			if index >= 10 {
				break
			}
			from := point
			to := point
			if index == 0 {
				from = current.Player.Position
			} else {
				from = current.GuidancePath[index-1]
			}
			if index+1 < len(current.GuidancePath) {
				to = current.GuidancePath[index+1]
			} else {
				to = point
			}
			guidance[point] = directionArrow(from, to)
		}
	}

	for y := 0; y < viewHeight; y++ {
		for x := 0; x < viewWidth; x++ {
			world := game.Point{X: worldLeft + x, Y: worldTop + y}
			tile := current.Maze.TileAt(world)
			character, style := tileAppearance(tile)
			if arrow, ok := guidance[world]; ok {
				character = arrow
				style = tcell.StyleDefault.Foreground(tcell.ColorAqua).Bold(true)
			}
			if world == current.Player.Position {
				character = '@'
				style = tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true)
			}
			screen.SetContent(viewLeft+x, viewTop+y, character, nil, style)
		}
	}

	DrawHUD(screen, current)
}

// tileAppearance maps a maze tile to its display character and style.
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
	case game.Food:
		return '%', tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true)
	default:
		return '█', tcell.StyleDefault.Foreground(tcell.ColorGray)
	}
}

// directionArrow points from one adjacent path cell toward the next.
func directionArrow(from, to game.Point) rune {
	switch {
	case to.X > from.X:
		return '→'
	case to.X < from.X:
		return '←'
	case to.Y > from.Y:
		return '↓'
	default:
		return '↑'
	}
}

// drawText writes a clipped line of text starting at the given screen position.
func drawText(screen tcell.Screen, x, y, screenWidth int, text string, style tcell.Style) {
	_, screenHeight := screen.Size()
	if y < 0 || y >= screenHeight {
		return
	}
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
