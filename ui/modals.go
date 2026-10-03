package ui

// This file draws title, confirmation, victory, and defeat overlays.

import (
	"fmt"
	"time"
	"unicode/utf8"

	"termzogs/game"

	"github.com/gdamore/tcell/v2"
)

// DrawOverlay renders title, item-confirmation, victory, or defeat overlays.
// It returns true when the overlay replaces the normal game view.
func DrawOverlay(screen tcell.Screen, current *game.Game) bool {
	return drawOverlayAt(screen, current, time.Now())
}

func drawOverlayAt(screen tcell.Screen, current *game.Game, now time.Time) bool {
	width, height := screen.Size()
	if width <= 0 || height <= 0 || current == nil {
		return true
	}

	var lines []string
	style := blackBackgroundStyle().Foreground(tcell.ColorWhite).Bold(true)
	switch current.State {
	case game.StateTitle:
		lines = append(termzogsTitleArt(), "©2026 By Kev Ellis - www.kevssite.com", "v1.0.0", "", "Find the treasure and return to the entrance.", "", "SPACE START  |  I INSTRUCTIONS  |  Q QUIT")
	case game.StateInstructions:
		lines = instructionsPageLines()
	case game.StateVictory:
		style = blackBackgroundStyle().Foreground(tcell.ColorGreen).Bold(true)
		lines = []string{
			"VICTORY! YOU ESCAPED WITH THE TREASURE!",
			fmt.Sprintf("Steps Taken: %d", current.Player.Moves),
			fmt.Sprintf("Zogs Slain: %d", current.ZogsSlain),
			fmt.Sprintf("Remaining Energy: %d", current.Player.Energy),
			"Press R to play again  |  Q to quit",
		}
	case game.StateDefeat:
		style = blackBackgroundStyle().Foreground(tcell.ColorRed).Bold(true)
		lines = []string{
			defeatHeadline(current.DefeatCause),
			fmt.Sprintf("Steps Taken: %d  |  Zogs Slain: %d", current.Player.Moves, current.ZogsSlain),
			"R: restart  |  M: main menu  |  Q: quit",
		}
	case game.StatePlaying:
		if !current.PendingSwap {
			return false
		}
		style = blackBackgroundStyle().Foreground(tcell.ColorYellow).Bold(true)
		lines = []string{"Drop Treasure to pick up Sword?", "Y: drop treasure and take sword  |  N: keep treasure"}
	default:
		return false
	}

	if current.State == game.StatePlaying {
		startY := height/2 - len(lines)/2
		for index, line := range lines {
			drawCentered(screen, width, startY+index, line, style)
		}
		return false
	}
	for index, line := range lines {
		lineStyle := style
		y := height/2 - len(lines)/2 + index
		if current.State == game.StateInstructions {
			switch index {
			case 0:
				lineStyle = blackBackgroundStyle().Foreground(tcell.ColorRed).Bold(true)
			case 2, 3, 4, 5, 6, 7, 8, 13, 14:
				lineStyle = blackBackgroundStyle().Foreground(tcell.ColorYellow).Bold(true)
			}
		} else if current.State == game.StateTitle && index < titleGlyphHeight {
			lineStyle = blackBackgroundStyle().Foreground(tcell.ColorRed).Bold(true)
		} else if current.State == game.StateTitle && index >= titleGlyphHeight && index <= titleGlyphHeight+1 {
			lineStyle = blackBackgroundStyle().Foreground(tcell.ColorBlue).Bold(true)
		}
		if current.State == game.StateTitle && index <= titleGlyphHeight+1 {
			y--
		}
		drawCentered(screen, width, y, line, lineStyle)
	}
	if current.State == game.StateTitle && screenTooSmall(width, height) && screenSizeWarningVisible(now) {
		warningStyle := blackBackgroundStyle().Foreground(tcell.ColorRed).Bold(true)
		drawCentered(screen, width, height-2, "SCREEN TOO SMALL", warningStyle)
		drawCentered(screen, width, height-1, "FULL VIEW: 60x20", warningStyle)
	}
	return true
}

func screenTooSmall(width, height int) bool {
	return width < MinimumScreenWidth || height < MinimumScreenHeight
}

func screenSizeWarningVisible(now time.Time) bool {
	return (now.UnixNano()/int64(screenSizeWarningFlashInterval))%2 == 0
}

const (
	titleGlyphWidth                = 6
	titleGlyphHeight               = 5
	titleWord                      = "TERMZOGS"
	titleArtWidth                  = len(titleWord)*titleGlyphWidth + len(titleWord) - 1
	screenSizeWarningFlashInterval = 500 * time.Millisecond
)

var titleGlyphs = map[rune][titleGlyphHeight]string{
	'T': {"TTTTTT", "  TT  ", "  TT  ", "  TT  ", "  TT  "},
	'E': {"EEEEEE", "EE    ", "EEEEE ", "EE    ", "EEEEEE"},
	'R': {"RRRRR ", "RR  RR", "RRRRR ", "RR RR ", "RR  RR"},
	'M': {"MM  MM", "MMMMMM", "MM  MM", "MM  MM", "MM  MM"},
	'Z': {"ZZZZZZ", "   ZZ ", "  ZZ  ", " ZZ   ", "ZZZZZZ"},
	'O': {" OOOO ", "OO  OO", "OO  OO", "OO  OO", " OOOO "},
	'G': {" GGGG ", "GG    ", "GG GGG", "GG  GG", " GGGG "},
	'S': {" SSSS ", "SS    ", " SSSS ", "    SS", "SSSSS "},
}

func termzogsTitleArt() []string {
	lines := make([]string, titleGlyphHeight)
	for _, letter := range titleWord {
		for row, glyph := range titleGlyphs[letter] {
			if lines[row] != "" {
				lines[row] += " "
			}
			lines[row] += glyph
		}
	}
	return lines
}

func instructionsPageLines() []string {
	return []string{
		"HOW TO PLAY TERMZOGS",
		"",
		"Find the treasure near the maze center and return to the entrance.",
		"Each move costs 1 energy. Food restores 50 energy.",
		"Armed: defeat Zogs with a sword; it breaks and combat costs 5 more.",
		"Unarmed fights are 50/50: victory costs 40 energy; defeat ends the run.",
		"You can carry a sword or treasure, but not both.",
		"Face a prisoner in a wall niche to ask for directions (cost: 15 energy).",
		"Guidance points to treasure, or to the entrance while carrying it.",
		"",
		"ARROWS / WASD Move     V/M Expanded map     Q/ESC Menu while playing",
		"SHIFT-D Drop item     Y/N Answer sword-swap prompt",
		"",
		"Termzogs is heavily inspired by the superb Sinclair ZX81 game",
		"Mazogs by Don Priestley.",
		"",
		"Press any key to return to title screen",
	}
}

func defeatHeadline(cause game.DefeatCause) string {
	if cause == game.DefeatCauseZog {
		return "GAME OVER - YOU WERE KILLED BY A ZOG"
	}
	return "GAME OVER - YOU PERISHED IN THE MAZE"
}

// drawCentered places a clipped text line at the horizontal screen midpoint.
func drawCentered(screen tcell.Screen, width, y int, text string, style tcell.Style) {
	drawText(screen, (width-utf8.RuneCountInString(text))/2, y, width, text, style)
}
