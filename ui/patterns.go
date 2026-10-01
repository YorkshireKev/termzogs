package ui

import "strings"

type cellPattern [ViewportCellHeight]string

// Each character's frames are ordered from first to last animation pose.
var localCellPatterns = map[rune][]cellPattern{
	// Floor
	' ': {
		{
			"      ",
			"      ",
			"      ",
		},
	},
	// Sword pickup
	'/': {
		{
			"  |   ",
			"--+==>",
			"  |   ",
		},
	},
	// Prisoner
	'P': {
		{
			"PPPPPP",
			"PPPPPP",
			"PPPPPP",
		},
	},
	// Zog frames alternate its feet.
	'M': {
		{
			" M  M ",
			"MMMMMM",
			"M MM M",
		},
		{
			" M  M ",
			"MMMMMM",
			"MM  MM",
		},
	},
	// Treasure pulses between its full and small shapes.
	'$': {
		{
			"$$$$$$",
			"$$$$$$",
			"$$$$$$",
		},
		{
			"      ",
			"  $$  ",
			"      ",
		},
		{
			" $$$$ ",
			" $$$$ ",
			"      ",
		},
	},
	// Food
	'%': {
		{
			"  /|  ",
			" /%%\\ ",
			" \\%%/ ",
		},
	},
	// Wall
	'█': {
		{
			"██████",
			"██████",
			"██████",
		},
	},
	// Player
	'@': {
		{
			"  O   ",
			" /|\\  ",
			" / \\  ",
		},
		{
			"  O   ",
			" /|\\  ",
			" | |  ",
		},
	},
	// Player facing left
	'<': {
		{
			"  <O  ",
			" <=|\\ ",
			" <  \\ ",
		},
		{
			"  <O  ",
			" <=|\\ ",
			" < /  ",
		},
	},
	// Player facing right
	'>': {
		{
			"  O>  ",
			" /|=> ",
			" /  > ",
		},
		{
			"  O>  ",
			" /|=> ",
			"  / > ",
		},
	},
	// Guidance arrow pointing right
	'→': {
		{
			"→→→→→→",
			"→→→→→→",
			"→→→→→→",
		},
	},
	// Guidance arrow pointing left
	'←': {
		{
			"←←←←←←",
			"←←←←←←",
			"←←←←←←",
		},
	},
	// Guidance arrow pointing down
	'↓': {
		{
			"↓↓↓↓↓↓",
			"↓↓↓↓↓↓",
			"↓↓↓↓↓↓",
		},
	},
	// Guidance arrow pointing up
	'↑': {
		{
			"↑↑↑↑↑↑",
			"↑↑↑↑↑↑",
			"↑↑↑↑↑↑",
		},
	},
}

func localCellPattern(character rune) cellPattern {
	return localCellPatternForFrame(character, 0)
}

func localCellPatternForFrame(character rune, frame int) cellPattern {
	if frames, ok := localCellPatterns[character]; ok && len(frames) > 0 {
		return frames[frame%len(frames)]
	}
	row := strings.Repeat(string(character), ViewportCellWidth)
	return cellPattern{row, row, row}
}
