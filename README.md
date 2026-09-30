# TermZogs

TermZogs is a terminal maze crawler written in Go with `tcell`. Explore the maze, find the treasure, and return to the entrance while managing energy and avoiding Zogs.

## Requirements

- Go 1.20 or later
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

## Controls

| Key | Action |
| --- | --- |
| Arrow keys or lowercase WASD (`w`, `a`, `s`, `d`) | Move one cell |
| V | Open the expanded maze view |
| Uppercase D (Shift+D) | Drop the held sword or treasure on the current cell |
| Y / N | Accept or decline dropping treasure to take a sword |
| Q or Escape | Quit |
| R | Restart after victory or defeat |

To ask a prisoner for directions, stand next to it and press the direction key toward its cell. Prisoners are in wall niches, so the player stays in place and cannot walk onto them. This costs 15 energy and displays up to 10 directional steps toward the treasure, or toward the entrance if carrying the treasure. Guidance remains visible for up to 15 movement steps or prisoner interactions.

The expanded view is 70 columns by 20 rows, clipped to the terminal size. Opening it costs 20 energy credits and requires more than 30 energy. It stays open until the next key press; that key also performs its normal action. Pressing V while the expanded view is open returns to the local view without another charge.

## Game Rules

- The generated maze is 62 columns by 42 rows. It contains one treasure near the center, 12 swords, 35 prisoners in wall niches, 18 Zogs, and 8-12 food items in dead ends.
- The local gameplay view is 10 by 10 cells and follows the player.
- Players start with 200 energy and can hold up to 300. Each successful movement costs 1 energy. Reaching 0 energy causes defeat.
- Collecting food consumes it and restores 50 energy, up to the 300-energy maximum.
- The player can carry a sword or the treasure, but not both. Walking over a sword while carrying the treasure prompts for confirmation. Y drops the treasure and takes the sword; N keeps the treasure. A sword encountered while already carrying a sword remains in the maze.
- Dropping an item leaves it on the player's current cell. A sword prevents treasure pickup until it is dropped.
- Fighting a Zog while armed defeats it, breaks the sword, and costs 5 additional energy. Without a sword, combat is a 50/50 roll: a win costs 40 additional energy; a loss causes immediate defeat. The movement into the Zog's cell also costs 1 energy.
- Pressing toward an adjacent prisoner triggers a shortest-path BFS hint. Hints point to the treasure before it is collected and to the entrance afterward.
- Returning to the entrance at `(1, 1)` while carrying the treasure wins the game. Zero energy or losing unarmed combat ends the run.
- Status messages disappear after six successful movement steps. Victory and defeat screens remain until the player restarts or quits.