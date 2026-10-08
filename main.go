// Package main starts TermZogs and manages terminal input and rendering.
package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"termzogs/game"
	"termzogs/ui"

	"github.com/gdamore/tcell/v2"
)

// main reports startup or runtime errors to stderr.
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run owns the terminal screen and processes input until the user exits.
func run() error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return fmt.Errorf("create terminal screen: %w", err)
	}
	if err := screen.Init(); err != nil {
		return fmt.Errorf("initialize terminal screen: %w", err)
	}
	defer screen.Fini()
	screen.DisableMouse()

	screen.SetTitle("TermZogs")
	current := game.NewTitleGame()
	renderer := ui.Renderer{}
	stopEvents := make(chan struct{})
	defer close(stopEvents)
	events := make(chan tcell.Event, 16)
	// Keep event polling independent so initial rendering is not blocked.
	go screen.ChannelEvents(events, stopEvents)
	ticker := time.NewTicker(time.Second / 30)
	defer ticker.Stop()

	draw := func() {
		renderer.Draw(screen, current)
		screen.Show()
	}
	draw()

	for {
		select {
		case event, ok := <-events:
			if !ok {
				return errors.New("terminal event stream closed")
			}
			switch event := event.(type) {
			case *tcell.EventKey:
				if handleKey(current, event) {
					return nil
				}
				draw()
			case *tcell.EventResize:
				screen.Sync()
				draw()
			case *tcell.EventInterrupt:
				return nil
			}
		case <-ticker.C:
			now := time.Now()
			current.ResolveCombatIfReady(now)
			current.UpdateZogs(now)
			draw()
		}
	}
}

// handleKey maps quit and movement keys to game actions.
func handleKey(current *game.Game, event *tcell.EventKey) bool {
	if current.State == game.StateInstructions {
		current.State = game.StateTitle
		return false
	}
	if current.ExpandedView {
		current.ExpandedView = false
		current.SetStatusMessage("")
		if event.Rune() == 'v' || event.Rune() == 'V' || event.Rune() == 'm' || event.Rune() == 'M' {
			return false
		}
	}
	if event.Key() == tcell.KeyCtrlC {
		return true
	}
	if current.State != game.StateTitle && (event.Rune() == 'q' || event.Rune() == 'Q' || event.Key() == tcell.KeyEscape) {
		current.State = game.StateTitle
		return false
	}
	if current.State == game.StateTitle && (event.Rune() == 'q' || event.Rune() == 'Q' || event.Key() == tcell.KeyEscape) {
		return true
	}
	if current.CombatInProgress() {
		return false
	}
	if current.State == game.StateTitle {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case ' ':
				current.Start()
			case 'i', 'I':
				current.State = game.StateInstructions
			}
		}
		return false
	}
	if current.State == game.StateDefeat || current.State == game.StateVictory {
		switch event.Rune() {
		case 'r', 'R':
			current.Restart()
		case 'm', 'M':
			if current.State == game.StateDefeat {
				current.State = game.StateTitle
			}
		}
		return false
	}
	if current.PendingSwap {
		switch event.Rune() {
		case 'y', 'Y':
			current.AnswerSwordSwap(true)
		case 'n', 'N':
			current.AnswerSwordSwap(false)
		}
		return false
	}

	switch event.Rune() {
	case 'v', 'V', 'm', 'M':
		current.ToggleMazeView()
		return false
	case 'D':
		current.DropItem()
		return false
	}

	dx, dy := 0, 0
	switch event.Key() {
	case tcell.KeyUp:
		dy = -1
	case tcell.KeyDown:
		dy = 1
	case tcell.KeyLeft:
		dx = -1
	case tcell.KeyRight:
		dx = 1
	default:
		switch event.Rune() {
		case 'w', 'W':
			dy = -1
		case 's', 'S':
			dy = 1
		case 'a', 'A':
			dx = -1
		case 'd':
			dx = 1
		}
	}
	if dx != 0 || dy != 0 {
		current.Move(dx, dy)
	}
	return false
}
