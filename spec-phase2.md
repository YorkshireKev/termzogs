# Product Specification Extension: TermZogs (Phase 2 - Mechanics & Systems)

## Feature Overview
This extension updates the core `TermZogs` codebase to incorporate full game mechanics based on the original 1982 *Mazogs* rules: player energy/food management, sword mechanics, monster combat, prisoner pathfinding guidance, and victory/defeat game states.

---

## 1. Updated Architecture & File Extensions

```
Expand the codebase with these additions:
termzogs/
├── game/
│   ├── combat.go       # Unarmed and armed combat resolution
│   ├── energy.go       # Food, step counters, and energy decay
│   └── pathfinding.go  # BFS algorithm for prisoner guidance routes
└── ui/
├── hud.go          # Status bar, energy meter, inventory indicators
└── modals.go       # Combat logs, victory/defeat dialog screens
```
## 2. Expanded Game State Rules

### A. Energy & Step System (`game/energy.go`)
- **Initial Energy:** Player starts with **200 Units** (Max capacity: 300).
- **Movement Cost:** Every movement step costs **1 Energy Unit**.
- **Food Placement:** 
  - Scatter 8–12 Food items (`%`) in maze dead-ends.
  - Picking up Food restores **+50 Energy** (capped at Max capacity).
- **Zero Energy State:** When Energy reaches 0, player dies immediately (**Defeat State**).

### B. Inventory & Item Interaction (`game/state.go`)
- **Exclusive Hold Constraint:** Player can hold **EITHER** a Sword (`/`) **OR** the Treasure Chest (`$`).
  - Attempting to pick up a Sword while holding the Treasure prompts: *"Drop Treasure to pick up Sword? [Y/N]"*.
  - Dropping an item leaves it on the current tile cell.

### C. Combat Resolution (`game/combat.go`)
Triggered whenever the player enters a cell occupied by a Zog (`M`):

1. **Armed Combat (Player holds Sword `/`):**
   - **Outcome:** Zog is removed from the maze grid.
   - **Cost:** Sword breaks (`HasSword = false`). Player loses **5 Energy** from exertion.
   - **UI Output:** Status message: *"You slew the Zog, but your sword shattered!"*

2. **Unarmed Combat (Player holds no Sword):**
   - **Mechanic:** 50/50 probability roll (`rand.Float32() < 0.5`).
   - **Victory (50% Chance):** Player defeats the Zog, but suffers a severe **-40 Energy** penalty.
   - **Defeat (50% Chance):** Player dies immediately (**Defeat State**).
   - **UI Output:** Status message: *"Desperate fight! You defeated the Zog but lost 40 Energy!"* or *"The Zog overpowered you!"*

### D. Prisoner Guidance System (`game/pathfinding.go`)
When the player moves adjacent to a Prisoner (`P`) and presses `E` (Interact):

1. **Cost:** Interacting costs **15 Energy**.
2. **Pathfinding Logic:**
   - If `HasTreasure == false`: Run **Breadth-First Search (BFS)** from player position to Treasure (`$`) coordinates.
   - If `HasTreasure == true`: Run **BFS** from player position to Start/Exit (`(1,1)`).
3. **Rendering Guidance:**
   - Highlight the next **10 steps** of the calculated path on the viewport using directional arrows (`↑`, `↓`, `←`, `→`) in bright cyan for **15 game ticks**.

---

## 3. UI & State Updates (`ui/hud.go` & `ui/modals.go`)

### A. Enhanced HUD Status Bar
Divide the bottom of the screen into three clear data fields:
[ Energy: ████████░░ 180/300 ]  [ Holding: SWORD ]  [ Treasure: NO ]


### B. Game State Transitions
1. **Title Screen State:** Press `Space` to start game generation.
2. **Playing State:** Active local viewport rendering and event loop.
3. **Victory State:** Player reaches cell `(1,1)` while `HasTreasure == true`.
   - Render ASCII overlay banner: **"VICTORY! YOU ESCAPED WITH THE TREASURE!"**
   - Display stats: Total Steps Taken, Zogs Slain, Remaining Energy.
4. **Defeat State:** Player dies from Combat or Zero Energy.
   - Render ASCII overlay banner: **"GAME OVER - YOU PERISHED IN THE MAZE"**
   - Prompt: Press `R` to restart or `Q` to quit.

---

## 4. Instructions for AI Code Agent

1. Update `game/state.go` with energy fields (`Energy`, `MaxEnergy`), item flags, and game status enums (`StateTitle`, `StatePlaying`, `StateVictory`, `StateDefeat`).
2. Build `game/energy.go` and `game/combat.go` handling move deductions and probability rolls.
3. Build `game/pathfinding.go` implementing BFS on the maze array to return an ordered slice of step coordinates.
4. Update `ui/renderer.go` and `ui/hud.go` to render the updated HUD bar, food icons (`%`), and path overlay arrows.
5. Wire key triggers: `E` for Prisoner interaction, `D` for dropping current item, `R` for restart on Game Over.