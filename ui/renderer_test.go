package ui

import (
	"strings"
	"testing"
	"time"

	"termzogs/game"

	"github.com/gdamore/tcell/v2"
)

func TestLocalCellPatternsMatchViewportCellSize(t *testing.T) {
	for character, frames := range localCellPatterns {
		if len(frames) == 0 {
			t.Errorf("pattern %q has no frames", character)
		}
		for frameIndex, pattern := range frames {
			for rowIndex, row := range pattern {
				if width := len([]rune(row)); width != ViewportCellWidth {
					t.Errorf("pattern %q frame %d row %d has width %d, want %d", character, frameIndex, rowIndex, width, ViewportCellWidth)
				}
			}
		}
	}
}

func TestTimedAnimationFramesAdvanceEachInterval(t *testing.T) {
	for frame := 0; frame < 4; frame++ {
		at := time.Duration(frame) * monsterFrameDuration
		if got := timedAnimationFrame(time.Unix(0, int64(at))); got != frame {
			t.Errorf("animation frame at %s = %d, want %d", at, got, frame)
		}
	}
}

func TestTreasurePulsesThroughThreeTimedFrames(t *testing.T) {
	fullSize := cellPattern{"$$$$$$", "$$$$$$", "$$$$$$"}
	smallSize := cellPattern{"      ", "  $$  ", "      "}
	mediumSize := cellPattern{" $$$$ ", " $$$$ ", "      "}
	if fullSize == smallSize || fullSize == mediumSize || smallSize == mediumSize {
		t.Fatal("treasure animation should have three distinct sizes")
	}
	want := []cellPattern{fullSize, smallSize, mediumSize, fullSize}
	for tick, expected := range want {
		at := time.Duration(tick) * monsterFrameDuration
		frame := timedAnimationFrame(time.Unix(0, int64(at)))
		if got := localCellPatternForFrame('$', frame); got != expected {
			t.Errorf("treasure frame at tick %d = %v, want %v", tick, got, expected)
		}
	}

	firstMonsterFrame := localCellPatternForFrame('M', 0)
	secondMonsterFrame := localCellPatternForFrame('M', 1)
	thirdMonsterFrame := localCellPatternForFrame('M', 2)
	if firstMonsterFrame != thirdMonsterFrame || firstMonsterFrame == secondMonsterFrame {
		t.Fatal("monster animation should continue looping through two poses")
	}
}

func TestPlayerAnimationsHaveTwoFramesAndMatchFacing(t *testing.T) {
	for _, sprite := range []rune{'<', '>', '@'} {
		if frames := localCellPatterns[sprite]; len(frames) != 2 {
			t.Errorf("player sprite %q has %d frames, want 2", sprite, len(frames))
		}
	}

	for _, test := range []struct {
		facing game.Point
		want   rune
	}{
		{facing: game.Point{X: -1}, want: '<'},
		{facing: game.Point{X: 1}, want: '>'},
		{facing: game.Point{Y: -1}, want: '@'},
		{facing: game.Point{Y: 1}, want: '@'},
	} {
		if got := playerSpriteForFacing(test.facing); got != test.want {
			t.Errorf("sprite for facing %v = %q, want %q", test.facing, got, test.want)
		}
	}
}

func TestCombatAnimationAlternatesPlayerAndZog(t *testing.T) {
	facing := game.Point{X: -1}
	for frame, want := range []rune{'M', '<', 'M', '<'} {
		if got := combatSpriteForFrame(facing, frame); got != want {
			t.Errorf("combat sprite at frame %d = %q, want %q", frame, got, want)
		}
	}
}

func TestDefeatHeadlineNamesZogCause(t *testing.T) {
	if got := defeatHeadline(game.DefeatCauseZog); got != "GAME OVER - YOU WERE KILLED BY A ZOG" {
		t.Errorf("Zog defeat headline = %q", got)
	}
	if got := defeatHeadline(game.DefeatCauseEnergy); got != "GAME OVER - YOU PERISHED IN THE MAZE" {
		t.Errorf("energy defeat headline = %q", got)
	}
}

func TestTermzogsTitleArtUsesSixColumnLetterGlyphs(t *testing.T) {
	lines := termzogsTitleArt()
	if len(lines) != titleGlyphHeight {
		t.Fatalf("title art has %d rows, want %d", len(lines), titleGlyphHeight)
	}
	for rowIndex, line := range lines {
		if width := len([]rune(line)); width != titleArtWidth {
			t.Errorf("title art row %d has width %d, want %d", rowIndex, width, titleArtWidth)
		}
	}
	for letter, glyph := range titleGlyphs {
		for rowIndex, line := range glyph {
			if width := len([]rune(line)); width != titleGlyphWidth {
				t.Errorf("glyph %q row %d has width %d, want %d", letter, rowIndex, width, titleGlyphWidth)
			}
			for _, character := range line {
				if character != letter && character != ' ' {
					t.Errorf("glyph %q contains unexpected character %q", letter, character)
				}
			}
		}
	}
}

func TestInstructionsPageIncludesControlsAndAttribution(t *testing.T) {
	lines := instructionsPageLines()
	page := strings.Join(lines, "\n")
	for _, expected := range []string{
		"ARROWS / WASD",
		"SHIFT-D",
		"V/M Expanded map",
		"Y/N",
		"Mazogs by Don Priestley",
		"Press I or ESC to return to title",
	} {
		if !strings.Contains(page, expected) {
			t.Errorf("instructions page is missing %q", expected)
		}
	}
	for row, line := range lines {
		if width := len([]rune(line)); width > 80 {
			t.Errorf("instructions row %d is %d columns wide, exceeds 80", row, width)
		}
	}
}

func TestTitleLogoIsRedAndSupportingTextIsWhite(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	if err := screen.Init(); err != nil {
		t.Fatal(err)
	}
	defer screen.Fini()

	Renderer{}.Draw(screen, &game.Game{State: game.StateTitle})
	screen.Show()
	cells, width, height := screen.GetContents()
	titleLines := append(termzogsTitleArt(), "©2026 By Kev Ellis - www.kevssite.com", "", "Find the treasure and return to the entrance.", "", "SPACE START  |  I INSTRUCTIONS  |  Q QUIT")
	startY := height/2 - len(titleLines)/2
	logoX := (width - titleArtWidth) / 2
	logoCell := cells[(startY-1)*width+logoX]
	logoForeground, _, _ := logoCell.Style.Decompose()
	if len(logoCell.Runes) == 0 || logoCell.Runes[0] != 'T' || logoForeground != tcell.ColorRed {
		t.Errorf("title logo cell = %q with foreground %v, want red T", logoCell.Runes, logoForeground)
	}

	copyright := titleLines[titleGlyphHeight]
	copyrightX := (width - len([]rune(copyright))) / 2
	copyrightCell := cells[(startY+titleGlyphHeight-1)*width+copyrightX]
	copyrightForeground, _, _ := copyrightCell.Style.Decompose()
	if len(copyrightCell.Runes) == 0 || copyrightCell.Runes[0] != '©' || copyrightForeground != tcell.ColorBlue {
		t.Errorf("copyright cell = %q with foreground %v, want blue ©", copyrightCell.Runes, copyrightForeground)
	}

	subtitleY := startY + titleGlyphHeight + 2
	subtitleX := (width - len(titleLines[titleGlyphHeight+2])) / 2
	subtitleCell := cells[subtitleY*width+subtitleX]
	subtitleForeground, _, _ := subtitleCell.Style.Decompose()
	if len(subtitleCell.Runes) == 0 || subtitleCell.Runes[0] != 'F' || subtitleForeground != tcell.ColorWhite {
		t.Errorf("subtitle cell = %q with foreground %v, want white F", subtitleCell.Runes, subtitleForeground)
	}
}

func TestInstructionsPageUsesRequestedTextColors(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	if err := screen.Init(); err != nil {
		t.Fatal(err)
	}
	defer screen.Fini()

	Renderer{}.Draw(screen, &game.Game{State: game.StateInstructions})
	screen.Show()
	cells, width, height := screen.GetContents()
	lines := instructionsPageLines()
	startY := height/2 - len(lines)/2

	for _, test := range []struct {
		line       int
		foreground tcell.Color
	}{
		{line: 0, foreground: tcell.ColorRed},
		{line: 2, foreground: tcell.ColorYellow},
		{line: 8, foreground: tcell.ColorYellow},
		{line: 10, foreground: tcell.ColorWhite},
		{line: 11, foreground: tcell.ColorWhite},
		{line: 13, foreground: tcell.ColorYellow},
		{line: 14, foreground: tcell.ColorYellow},
	} {
		x := (width - len([]rune(lines[test.line]))) / 2
		cell := cells[(startY+test.line)*width+x]
		foreground, _, _ := cell.Style.Decompose()
		if foreground != test.foreground {
			t.Errorf("line %q has foreground %v, want %v", lines[test.line], foreground, test.foreground)
		}
	}
}

func TestDrawAdvancesPlayerFrameWithMoveCount(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	if err := screen.Init(); err != nil {
		t.Fatal(err)
	}
	defer screen.Fini()

	current := &game.Game{
		Maze:  &game.Maze{},
		State: game.StatePlaying,
		Player: game.Player{
			Position: game.Point{X: 10, Y: 10},
		},
	}

	for moves := 0; moves < 3; moves++ {
		current.Player.Moves = moves
		Renderer{}.Draw(screen, current)
		screen.Show()
		cells, width, height := screen.GetContents()
		viewLeft := (width - ViewportWidth*ViewportCellWidth) / 2
		viewTop := (height - 2 - ViewportHeight*ViewportCellHeight) / 2
		worldLeft := current.Player.Position.X - ViewportWidth/2
		worldTop := current.Player.Position.Y - ViewportHeight/2
		cellLeft := viewLeft + (current.Player.Position.X-worldLeft)*ViewportCellWidth
		cellTop := viewTop + (current.Player.Position.Y-worldTop)*ViewportCellHeight
		expected := localCellPatterns['@'][moves%len(localCellPatterns['@'])]
		for row, expectedRow := range expected {
			for column, want := range []rune(expectedRow) {
				if got := cells[(cellTop+row)*width+cellLeft+column].Runes[0]; got != want {
					t.Errorf("moves=%d player pattern at (%d, %d) = %q, want %q", moves, column, row, got, want)
				}
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
	originalPatterns := localCellPatterns['M']
	customPattern := cellPattern{"ABCDEF", "GHIJKL", "MNOPQR"}
	localCellPatterns['M'] = []cellPattern{customPattern, customPattern}
	defer func() {
		localCellPatterns['M'] = originalPatterns
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
	for row, expected := range localCellPatterns['M'][0] {
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
