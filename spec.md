# Product Specification: TermZogs (TUI Maze Crawler in Go)

## Project Overview
TermZogs is a CLI terminal game written in Go using `github.com/gdamore/tcell/v2`. It is an open-source homage to the 1982 ZX81 classic "Mazogs". The player navigates a generated maze to find a hidden treasure chest and escape back to the starting point while dealing with wandering monsters ("Zogs"), swords, and prisoners.

---

## Technical Stack & Constraints
1. **Language:** Go (1.20+)
2. **TUI Library:** `github.com/gdamore/tcell/v2`
3. **Dependencies:** Pure Go (no CGO required).
4. **Target Interface:** Standard ANSI terminal with raw mode input and UTF-8 rendering support.

---

## Game Architecture & File Structure
Organize the codebase cleanly into the following structure:

```text
termzogs/
├── main.go               # Entry point, event loop, and tcell setup
├── game/
│   ├── game.go           # Core Game struct and main tick/update engine
│   ├── maze.go           # Maze generation algorithm (Recursive Backtracking)
│   ├── pathfinding.go    # BFS solver for prisoner guidance paths
│   └── state.go          # Player, Monster, Sword, and Item structs
├── ui/
│   ├── renderer.go       # Screen drawing functions using tcell SetContent
│   └── components.go     # Title screen, sidebar UI, status bar, and dialogs
└── go.mod
```

---

## Core Game Mechanics

### 1. The Maze Engine (`game/maze.go`)
- **Dimensions:** 62x62 cell grid surrounded by outer wall boundaries.
- **Generation:** Recursive backtracking creates the connected maze; randomized extra passages and small open rooms add alternate routes while preserving branches and dead ends.
- **Key Entities Placed at Generation:**
  - **Start Coordinate:** Edge cell at (1, 1).
  - **Treasure Chest (`$`):** Placed near the center of the maze.
  - **Swords (`/`):** Scattered throughout corridors (~10-15 per maze).
  - **Prisoners (`P`):** Trapped in wall niches (~5 per maze).
  - **Zogs (`M`):** Monsters placed in random corridors (~15-20 per maze).

### 2. Entity Systems & Rules (`game/state.go`)
- **Player (`@`):**
  - Moves via Arrow Keys or WASD.
  - Holds state flags: `HasSword` (bool), `HasTreasure` (bool), `Moves` (int step counter).
  - *Inventory Rule:* Player can ONLY hold a sword OR the treasure chest at one time (cannot hold both).
- **Zogs (`M`):**
  - Patrol adjacent open corridor cells on every Nth game tick.
  - *Combat Rule (With Sword):* Defeating a Zog consumes/breaks the sword (`HasSword = false`).
  - *Combat Rule (Without Sword):* Bumping into a Zog without a sword triggers a probability roll or combat interaction where the player risks death.
- **Prisoners (`P`):**
  - Interacting with a prisoner triggers a Breadth-First Search (BFS) pathfinding algorithm (`game/pathfinding.go`).
  - Highlights the shortest path toward the treasure (or toward the exit if the treasure is already collected) using directional arrows (`↑`, `↓`, `←`, `→`).

---

## User Interface & Rendering (`ui/renderer.go`)

### 1. Dual Camera Views
- **Zoomed Local View (Default Gameplay):** Render a 15x15 cell viewport centered on the player's current coordinates.
- **Overview Map (Key `M`):** Toggle a zoomed-out view showing the entire 62x62 maze layout, rendering walls as small dots (`.`) and the player as a bright marker.

### 2. TUI Cell Character Mapping

| Entity | Character | `tcell` Foreground Color |
|---|---|---|
| Solid Wall | `█` or `▓` | `ColorDarkGray` |
| Corridor | `' '` (Space) | `ColorBlack` |
| Player | `@` | `ColorYellow` (Bold) |
| Zog (Monster) | `M` | `ColorRed` (Bold) |
| Sword | `/` | `ColorCyan` |
| Prisoner | `P` | `ColorGreen` |
| Treasure | `$` | `ColorGold` / `ColorYellow` |
| Path Arrow | `↑`, `↓`, `←`, `→` | `ColorLightBlue` |

---

## Instructions for AI Agent

1. Begin by generating `go.mod` and setting up `main.go` with a clean `tcell` initialization/de-initialization pattern (`defer screen.Fini()`).
2. Implement the `Game` state struct and event loop with non-blocking key listening and ticker-based rendering (~30 FPS).
3. Build `game/maze.go` to generate a valid 62x62 maze array and place all items/entities.
4. Build the local viewport rendering engine in `ui/renderer.go` to draw walls and entities onto the screen.
5. Add movement, collision detection, and combat state handling.
