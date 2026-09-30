package ui

import "strings"

type cellPattern [ViewportCellHeight]string

// Edit these three six-character rows to draw local gameplay sprites.
var localCellPatterns = map[rune]cellPattern{
	// Floor
	' ': {
		"      ",
		"      ",
		"      ",
	},
	// Sword pickup
	'/': {
		"//////",
		"//////",
		"//////",
	},
	// Prisoner
	'P': {
		"PPPPPP",
		"PPPPPP",
		"PPPPPP",
	},
	// Zog
	'M': {
		"MMMMMM",
		"MMMMMM",
		"MMMMMM",
	},
	// Treasure
	'$': {
		"$$$$$$",
		"$$$$$$",
		"$$$$$$",
	},
	// Food
	'%': {
		"%%%%%%",
		"%%%%%%",
		"%%%%%%",
	},
	// Wall
	'█': {
		"██████",
		"██████",
		"██████",
	},
	// Player
	'@': {
		"  O   ",
		" /|\\  ",
		" / \\  ",
	},
	// Guidance arrow pointing right
	'→': {
		"→→→→→→",
		"→→→→→→",
		"→→→→→→",
	},
	// Guidance arrow pointing left
	'←': {
		"←←←←←←",
		"←←←←←←",
		"←←←←←←",
	},
	// Guidance arrow pointing down
	'↓': {
		"↓↓↓↓↓↓",
		"↓↓↓↓↓↓",
		"↓↓↓↓↓↓",
	},
	// Guidance arrow pointing up
	'↑': {
		"↑↑↑↑↑↑",
		"↑↑↑↑↑↑",
		"↑↑↑↑↑↑",
	},
}

func localCellPattern(character rune) cellPattern {
	if pattern, ok := localCellPatterns[character]; ok {
		return pattern
	}
	row := strings.Repeat(string(character), ViewportCellWidth)
	return cellPattern{row, row, row}
}
