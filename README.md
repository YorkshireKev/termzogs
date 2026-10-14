# TermZogs

Termzogs is a terminal maze crawler written in Go with `tcell`. Explore the maze, find the treasure, and return to the entrance while managing energy and avoiding Zogs. This was an experiment in Vibe Coding a terminal go based game. You can read more about how it was made on by blog at [kevssite.com/termzogs-game](https://kevssite.com/termzogs-game).

Termzogs is heavily influenced by the excellent ZX81 Classic game Mazogs, by Don Priestley. I should be clear, this is not a port of Mazogs, it is a much poorer game inspired by a masterpiece. I encourage you to give Mazogs a play on a ZX81 emulator or, even better, a real ZX81.

If you just want to see what the game looks and plays like, you can play it in a [web-browser](https://kevssite.com/termzogs).

## Requirements

- Go 1.23 or later
- An ANSI-compatible terminal with UTF-8 support

## Run

From the project directory:

```sh
go run .
```

At the title screen, press Space to generate and start a maze. To run the tests:

```sh
go test ./...
```

The full-size local game view requires a terminal at least 60 columns by 20 rows. A warning appears on the title screen if the terminal is smaller.

The web build scales the terminal text to fit smaller browser windows while keeping its 80-by-24 grid.

## Controls

| Key | Action |
| --- | --- |
| Arrow keys or lowercase WASD (`w`, `a`, `s`, `d`) | Move one cell |
| V or M | Open the expanded maze view |
| Uppercase D (Shift+D) | Drop the held sword or treasure on the current cell |
| Y / N | Accept or decline dropping treasure to take a sword |
| Q | Quit |
| Escape | Quit in terminal builds; disabled in the web version |
| R | Restart after victory or defeat |
| M | Return to the title screen after defeat |

To ask a prisoner for directions, stand next to it and press the direction key toward its cell. Prisoners are in wall niches, so the player stays in place and cannot walk onto them. This costs 15 energy and displays up to 10 directional steps toward the treasure, or toward the entrance if carrying the treasure. Guidance remains visible for up to 15 movement steps or prisoner interactions.

The expanded view is 70 columns by 20 rows, clipped to the terminal size. Opening it costs 20 energy credits and requires more than 30 energy. It stays open until the next key press; that key also performs its normal action. Pressing V while the expanded view is open returns to the local view without another charge.

## Game Rules

- The generated maze is 62 columns by 42 rows, with randomized shortcuts and open rooms connecting its branching corridors. It contains one treasure near the center, 18 swords, 35 prisoners in wall niches, and 18 Zogs. Two Zogs guard the main start-to-treasure route, each preceded by a sword; other Zogs and swords are placed off-route. It contains 15–20 food supplies, including enough on the main route to cover its armed fights.
- Zogs move one cell per second in repeating patrol loops of 6–10 steps, kept near their starting locations. They pause while the expanded map is open.
- The local gameplay view follows the player across 10 by 6 cells. Each cell is rendered as a 6 by 3 block of its tile character; the expanded map view remains unchanged.
- Players start with 250 energy and can hold up to 350 on typical mazes. Capacity scales up for unusually long routes to ensure the guaranteed supplies can cover the return trip. Each successful movement costs 1 energy. Reaching 0 energy causes defeat.
- Collecting food restores 50 energy, up to the current energy capacity.
- The player can carry a sword or the treasure, but not both. Walking over a sword while carrying the treasure prompts for confirmation. Y drops the treasure and takes the sword; N keeps the treasure. A sword encountered while already carrying a sword remains in the maze.
- Dropping an item leaves it on the player's current cell. Holding a sword prevents treasure pickup until it is dropped.
- Fighting a Zog while armed defeats it, but breaks the sword, and costs 5 additional energy. Without a sword, combat is a 50/50 roll: a win costs 40 additional energy; a loss causes immediate defeat.
- Prisoner point the way to the treasure before it is collected and to the entrance afterward. Asking for the directions hint costs 15 energy.
- Returning to the entrance at while carrying the treasure wins the game. Zero energy or losing unarmed combat ends the run.