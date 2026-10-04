package ui

// This file draws the player-centered maze viewport and status line.

import (
	"termzogs/game"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	ViewportWidth          = 10
	ViewportHeight         = 6
	ViewportCellWidth      = 6
	ViewportCellHeight     = 3
	ExpandedViewportWidth  = 70
	ExpandedViewportHeight = 20
	monsterFrameDuration   = 350 * time.Millisecond
	MinimumScreenWidth     = ViewportWidth * ViewportCellWidth
	MinimumScreenHeight    = ViewportHeight*ViewportCellHeight + 2
)

// Renderer draws game state onto a tcell screen.
type Renderer struct{}

// Draw paints the local viewport and player status.
func (Renderer) Draw(screen tcell.Screen, current *game.Game) {
	screenWidth, screenHeight := screen.Size()
	screen.SetStyle(blackBackgroundStyle())
	screen.Clear()
	if screenWidth <= 0 || screenHeight <= 0 {
		return
	}
	if DrawOverlay(screen, current) {
		return
	}

	if current.ExpandedView {
		drawExpandedViewport(screen, current)
	} else {
		drawLocalViewport(screen, current)
	}
	DrawHUD(screen, current)
}

func drawExpandedViewport(screen tcell.Screen, current *game.Game) {
	screenWidth, screenHeight := screen.Size()
	viewWidth := ExpandedViewportWidth
	viewHeight := ExpandedViewportHeight
	if viewWidth > screenWidth {
		viewWidth = screenWidth
	}
	if viewHeight > screenHeight {
		viewHeight = screenHeight
	}
	if viewHeight > screenHeight-3 {
		viewHeight = screenHeight - 3
	}
	if viewHeight <= 0 {
		return
	}
	viewLeft := (screenWidth - viewWidth) / 2
	viewTop := (screenHeight-viewHeight)/2 - 1
	if viewTop < 0 {
		viewTop = 0
	}
	worldLeft := current.Player.Position.X - viewWidth/2
	worldTop := current.Player.Position.Y - viewHeight/2
	guidance := guidanceFor(current)

	for y := 0; y < viewHeight; y++ {
		for x := 0; x < viewWidth; x++ {
			world := game.Point{X: worldLeft + x, Y: worldTop + y}
			character, style := displayAt(current, world, guidance)
			screen.SetContent(viewLeft+x, viewTop+y, character, nil, style)
		}
	}
}

func drawLocalViewport(screen tcell.Screen, current *game.Game) {
	screenWidth, screenHeight := screen.Size()
	viewWidth := ViewportWidth
	viewHeight := ViewportHeight
	if screenWidth/ViewportCellWidth < viewWidth {
		viewWidth = screenWidth / ViewportCellWidth
	}
	playAreaHeight := screenHeight - 2
	if playAreaHeight/ViewportCellHeight < viewHeight {
		viewHeight = playAreaHeight / ViewportCellHeight
	}
	if viewWidth <= 0 || viewHeight <= 0 {
		return
	}
	drawWidth := viewWidth * ViewportCellWidth
	drawHeight := viewHeight * ViewportCellHeight
	viewLeft := (screenWidth - drawWidth) / 2
	viewTop := (playAreaHeight - drawHeight) / 2
	worldLeft := current.Player.Position.X - viewWidth/2
	worldTop := current.Player.Position.Y - viewHeight/2
	guidance := guidanceFor(current)
	timedFrame := timedAnimationFrame(time.Now())

	for y := 0; y < viewHeight; y++ {
		for x := 0; x < viewWidth; x++ {
			world := game.Point{X: worldLeft + x, Y: worldTop + y}
			character, style := displayAt(current, world, guidance)
			frame := timedFrame
			if current.CombatInProgress() && world == current.Player.Position {
				character = combatSpriteForFrame(current.Player.Facing, timedFrame)
				if character == 'M' {
					style = blackBackgroundStyle().Foreground(tcell.ColorRed).Bold(true)
					frame = timedFrame
				} else {
					style = blackBackgroundStyle().Foreground(tcell.ColorYellow).Bold(true)
					frame = current.Player.Moves
				}
			} else if world == current.Player.Position {
				character = playerSpriteForFacing(current.Player.Facing)
				frame = current.Player.Moves
			}
			pattern := localCellPatternForFrame(character, frame)
			for blockY, row := range pattern {
				characters := []rune(row)
				for blockX := 0; blockX < ViewportCellWidth; blockX++ {
					patternCharacter := ' '
					if blockX < len(characters) {
						patternCharacter = characters[blockX]
					}
					screen.SetContent(viewLeft+x*ViewportCellWidth+blockX, viewTop+y*ViewportCellHeight+blockY, patternCharacter, nil, style)
				}
			}
		}
	}
}

func playerSpriteForFacing(facing game.Point) rune {
	switch {
	case facing.X < 0:
		return '<'
	case facing.X > 0:
		return '>'
	default:
		return '@'
	}
}

func combatSpriteForFrame(facing game.Point, frame int) rune {
	if frame%2 == 0 {
		return 'M'
	}
	return playerSpriteForFacing(facing)
}

func timedAnimationFrame(now time.Time) int {
	return int(now.UnixNano() / int64(monsterFrameDuration))
}

func guidanceFor(current *game.Game) map[game.Point]rune {
	guidance := make(map[game.Point]rune)
	if current.GuidanceTicks > 0 {
		for index, point := range current.GuidancePath {
			if index >= 10 {
				break
			}
			from, to := current.Player.Position, point
			if index+1 < len(current.GuidancePath) {
				from = point
				to = current.GuidancePath[index+1]
			} else if index > 0 {
				from = current.GuidancePath[index-1]
			} else {
				to = point
			}
			guidance[point] = directionArrow(from, to)
		}
	}
	return guidance
}

func displayAt(current *game.Game, world game.Point, guidance map[game.Point]rune) (rune, tcell.Style) {
	tile := current.Maze.TileAt(world)
	character, style := tileAppearance(tile)
	if arrow, ok := guidance[world]; ok && tile == game.Floor {
		character = arrow
		style = blackBackgroundStyle().Foreground(tcell.ColorAqua).Bold(true)
	}
	if world == current.Player.Position {
		character = '@'
		style = blackBackgroundStyle().Foreground(tcell.ColorYellow).Bold(true)
	}
	return character, style
}

// tileAppearance maps a maze tile to its display character and style.
func tileAppearance(tile game.Tile) (rune, tcell.Style) {
	switch tile {
	case game.Floor:
		return ' ', blackBackgroundStyle().Foreground(tcell.ColorBlack)
	case game.Sword:
		return '/', blackBackgroundStyle().Foreground(tcell.ColorAqua)
	case game.Prisoner:
		return 'P', blackBackgroundStyle().Foreground(tcell.ColorGreen)
	case game.Monster:
		return 'M', blackBackgroundStyle().Foreground(tcell.ColorRed).Bold(true)
	case game.Treasure:
		return '$', blackBackgroundStyle().Foreground(tcell.ColorYellow).Bold(true)
	case game.Food:
		return '%', blackBackgroundStyle().Foreground(tcell.ColorGreen).Bold(true)
	default:
		return '█', blackBackgroundStyle().Foreground(tcell.ColorGray)
	}
}

func blackBackgroundStyle() tcell.Style {
	return tcell.StyleDefault.Background(tcell.ColorBlack)
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
	if text == "" {
		if x >= 0 && x < screenWidth {
			screen.SetContent(x, y, ' ', nil, style)
		}
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
