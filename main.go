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
			draw()
		}
	}
}

// handleKey maps quit and movement keys to game actions.
func handleKey(current *game.Game, event *tcell.EventKey) bool {
	if event.Key() == tcell.KeyCtrlC || event.Key() == tcell.KeyEscape || event.Rune() == 'q' || event.Rune() == 'Q' {
		return true
	}
	if current.State == game.StateTitle {
		if event.Key() == tcell.KeyRune && event.Rune() == ' ' {
			current.Start()
		}
		return false
	}
	if current.State == game.StateDefeat || current.State == game.StateVictory {
		if event.Rune() == 'r' || event.Rune() == 'R' {
			current.Restart()
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
	case 'e', 'E':
		current.Interact()
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
