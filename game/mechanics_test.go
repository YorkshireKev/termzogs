package game

import (
	"math/rand"
	"testing"
	"time"
)

func TestMoveTracksFacingDirection(t *testing.T) {
	current := NewGameWithSeed(10)
	current.Player.Position = Point{X: 2, Y: 2}
	current.Maze.Tiles[2][2] = Floor
	directions := []Point{{X: 1}, {Y: 1}, {X: -1}, {Y: -1}}

	for _, direction := range directions {
		next := Point{X: current.Player.Position.X + direction.X, Y: current.Player.Position.Y + direction.Y}
		current.Maze.Tiles[next.Y][next.X] = Floor
		current.Move(direction.X, direction.Y)
		if current.Player.Facing != direction {
			t.Fatalf("facing = %v after move %v, want %v", current.Player.Facing, direction, direction)
		}
	}
}

func TestMoveEnergyFoodAndDefeat(t *testing.T) {
	current := NewGameWithSeed(11)
	current.Maze.Tiles[1][2] = Food
	current.Move(1, 0)
	if current.Player.Energy != InitialEnergy-1+50 {
		t.Fatalf("energy after food = %d, want %d", current.Player.Energy, InitialEnergy-1+50)
	}
	if current.Maze.TileAt(current.Player.Position) != Floor {
		t.Fatal("food tile was not consumed")
	}

	current.Player.Energy = 1
	current.Maze.Tiles[1][3] = Floor
	current.Move(1, 0)
	if current.State != StateDefeat || current.DefeatCause != DefeatCauseEnergy || current.Player.Energy != 0 {
		t.Fatalf("zero-energy move state = %v, cause = %v, energy = %d", current.State, current.DefeatCause, current.Player.Energy)
	}
}

func TestArmedCombatBreaksSwordAndRemovesZog(t *testing.T) {
	current := NewGameWithSeed(12)
	current.Player.HasSword = true
	current.Maze.Tiles[1][2] = Monster
	current.Move(1, 0)
	finishCombatForTest(t, current)
	if current.Player.HasSword || current.ZogsSlain != 1 || current.Maze.TileAt(current.Player.Position) != Floor {
		t.Fatal("armed combat did not remove the Zog and break the sword")
	}
	if current.Player.Energy != InitialEnergy-6 {
		t.Fatalf("energy after armed combat = %d, want %d", current.Player.Energy, InitialEnergy-6)
	}
}

func TestUnarmedCombatUsesSeededOutcome(t *testing.T) {
	current := NewGameWithSeed(13)
	current.rng = rand.New(rand.NewSource(1))
	current.Maze.Tiles[1][2] = Monster
	current.Move(1, 0)
	finishCombatForTest(t, current)
	if current.State != StateDefeat || current.DefeatCause != DefeatCauseZog || current.StatusMessage != "The Zog overpowered you!" {
		t.Fatalf("unarmed combat state = %v, cause = %v, message = %q", current.State, current.DefeatCause, current.StatusMessage)
	}
}

func TestUnarmedCombatCanWin(t *testing.T) {
	seed := int64(0)
	for ; seed < 100; seed++ {
		if rand.New(rand.NewSource(seed)).Float32() < 0.5 {
			break
		}
	}
	current := NewGameWithSeed(14)
	current.rng = rand.New(rand.NewSource(seed))
	current.Maze.Tiles[1][2] = Monster
	current.Move(1, 0)
	finishCombatForTest(t, current)
	if current.State != StatePlaying || current.ZogsSlain != 1 || current.Player.Energy != InitialEnergy-41 {
		t.Fatalf("unarmed win state = %v, kills = %d, energy = %d", current.State, current.ZogsSlain, current.Player.Energy)
	}
}

func TestCombatOutcomeWaitsForAnimation(t *testing.T) {
	current := NewGameWithSeed(25)
	current.Player.HasSword = true
	current.Maze.Tiles[1][2] = Monster
	current.Move(1, 0)
	startedAt := current.combatStartedAt
	if !current.CombatInProgress() || current.ZogsSlain != 0 {
		t.Fatal("combat should be pending without revealing its outcome")
	}

	current.Move(0, 1)
	if current.Player.Position != (Point{X: 2, Y: 1}) || current.Player.Moves != 1 {
		t.Fatal("player moved again during the fight animation")
	}
	if current.ResolveCombatIfReady(startedAt.Add(CombatAnimationDuration - time.Nanosecond)) {
		t.Fatal("combat resolved before the animation duration elapsed")
	}
	if !current.ResolveCombatIfReady(startedAt.Add(CombatAnimationDuration)) {
		t.Fatal("combat did not resolve when the animation duration elapsed")
	}
	if current.CombatInProgress() || current.ZogsSlain != 1 {
		t.Fatal("combat outcome was not revealed after the animation")
	}
}

func finishCombatForTest(t *testing.T, current *Game) {
	t.Helper()
	if !current.ResolveCombatIfReady(current.combatStartedAt.Add(CombatAnimationDuration)) {
		t.Fatal("combat did not resolve")
	}
}

func TestPrisonerInteractionCostsEnergyAndGuides(t *testing.T) {
	current := NewGameWithSeed(15)
	for y := range current.Maze.Tiles {
		for x := range current.Maze.Tiles[y] {
			current.Maze.Tiles[y][x] = Wall
		}
	}
	for y := 1; y <= 5; y++ {
		for x := 1; x <= 5; x++ {
			current.Maze.Tiles[y][x] = Floor
		}
	}
	current.Maze.Start = Point{X: 1, Y: 1}
	current.Maze.Treasure = Point{X: 5, Y: 5}
	current.Player.Position = Point{X: 2, Y: 2}
	current.Maze.Tiles[2][3] = Prisoner

	current.Move(1, 0)
	if current.Player.Energy != InitialEnergy-15 || current.GuidanceTicks != 15 || len(current.GuidancePath) == 0 {
		t.Fatalf("interaction energy=%d ticks=%d path length=%d", current.Player.Energy, current.GuidanceTicks, len(current.GuidancePath))
	}
	if current.Player.Position != (Point{X: 2, Y: 2}) {
		t.Fatalf("player moved into prisoner wall cell: %v", current.Player.Position)
	}
	if current.GuidancePath[len(current.GuidancePath)-1] != current.Maze.Treasure {
		t.Fatal("guidance path does not end at the treasure")
	}

	current.Player.HasTreasure = true
	current.Move(1, 0)
	if current.GuidancePath[len(current.GuidancePath)-1] != current.Maze.Start {
		t.Fatal("guidance path does not end at the exit while carrying treasure")
	}
}

func TestSwordSwapAndReturnToStartWins(t *testing.T) {
	current := NewGameWithSeed(16)
	current.Maze.Tiles[1][2] = Sword
	current.Player.HasTreasure = true
	current.Move(1, 0)
	if !current.PendingSwap {
		t.Fatal("picking up a sword while carrying treasure did not prompt")
	}
	current.AnswerSwordSwap(true)
	if current.PendingSwap || !current.Player.HasSword || current.Player.HasTreasure || current.Maze.TileAt(current.Player.Position) != Treasure {
		t.Fatal("accepting sword swap did not drop treasure and take sword")
	}

	current.Player.HasSword = false
	current.Player.HasTreasure = true
	current.Player.Position = Point{X: 2, Y: 1}
	current.Maze.Tiles[1][1] = Floor
	current.Move(-1, 0)
	if current.State != StateVictory {
		t.Fatalf("returning to start with treasure state = %v, want victory", current.State)
	}
}

func TestDropItemLeavesItOnCurrentCell(t *testing.T) {
	current := NewGameWithSeed(17)
	current.Player.Position = Point{X: 2, Y: 2}
	current.Maze.Tiles[2][2] = Floor
	current.Player.HasSword = true
	current.DropItem()
	if current.Player.HasSword || current.Maze.TileAt(current.Player.Position) != Sword {
		t.Fatal("dropping the sword did not leave it on the current cell")
	}
}

func TestExpandedMazeViewCostsCreditsAndEnforcesMinimum(t *testing.T) {
	current := NewGameWithSeed(19)
	current.ToggleMazeView()
	if !current.ExpandedView || current.Player.Energy != InitialEnergy-MazeViewCost {
		t.Fatalf("view=%t energy=%d, want expanded view and %d energy", current.ExpandedView, current.Player.Energy, InitialEnergy-MazeViewCost)
	}
	current.ToggleMazeView()
	if current.ExpandedView || current.Player.Energy != InitialEnergy-MazeViewCost {
		t.Fatal("closing expanded view should not charge credits")
	}

	current.Player.Energy = MazeViewMinCredits
	current.ToggleMazeView()
	if current.ExpandedView || current.Player.Energy != MazeViewMinCredits {
		t.Fatal("insufficient credits opened the view or reduced energy")
	}
	if current.StatusMessage != "You don't have enough health credits to view the maze." {
		t.Fatalf("insufficient-credit message = %q", current.StatusMessage)
	}
}

func TestExpandedMazeViewAllowsCreditsAboveMinimum(t *testing.T) {
	current := NewGameWithSeed(20)
	current.Player.Energy = MazeViewMinCredits + 1
	current.ToggleMazeView()
	if !current.ExpandedView || current.Player.Energy != MazeViewMinCredits+1-MazeViewCost {
		t.Fatalf("view=%t energy=%d, expected expanded view with %d credits", current.ExpandedView, current.Player.Energy, MazeViewMinCredits+1-MazeViewCost)
	}
}

func TestStatusMessageExpiresAfterSixSuccessfulMoves(t *testing.T) {
	current := NewGameWithSeed(23)
	for x := 2; x <= 7; x++ {
		current.Maze.Tiles[1][x] = Floor
	}
	current.SetStatusMessage("Temporary message")
	for step := 0; step < StatusMessageLifetime-1; step++ {
		current.Move(1, 0)
		if current.StatusMessage == "" {
			t.Fatalf("message disappeared after %d moves", step+1)
		}
	}
	current.Move(1, 0)
	if current.StatusMessage != "" {
		t.Fatalf("message remained after %d moves: %q", StatusMessageLifetime, current.StatusMessage)
	}
}

func TestNewStatusMessageResetsSixMoveLifetime(t *testing.T) {
	current := NewGameWithSeed(24)
	for x := 2; x <= 9; x++ {
		current.Maze.Tiles[1][x] = Floor
	}
	current.SetStatusMessage("First message")
	for step := 0; step < 3; step++ {
		current.Move(1, 0)
	}
	current.SetStatusMessage("Replacement message")
	for step := 0; step < StatusMessageLifetime-1; step++ {
		current.Move(1, 0)
		if current.StatusMessage == "" {
			t.Fatalf("replacement message disappeared after %d moves", step+1)
		}
	}
	current.Move(1, 0)
	if current.StatusMessage != "" {
		t.Fatalf("replacement message remained after six moves: %q", current.StatusMessage)
	}
}

func TestCarriedSwordDoesNotConsumeAnotherSword(t *testing.T) {
	current := NewGameWithSeed(18)
	current.Player.HasSword = true
	current.Maze.Tiles[1][2] = Sword
	current.Maze.Tiles[1][3] = Floor

	current.Move(1, 0)
	if !current.Player.HasSword || current.Maze.TileAt(Point{X: 2, Y: 1}) != Sword {
		t.Fatal("walking over another sword consumed it while already armed")
	}

	current.Move(1, 0)
	current.DropItem()
	current.Move(-1, 0)
	if !current.Player.HasSword || current.Maze.TileAt(Point{X: 2, Y: 1}) != Floor {
		t.Fatal("uncollected sword could not be picked up after dropping the carried sword")
	}
}

func TestFindPathReturnsOrderedShortestRoute(t *testing.T) {
	maze := &Maze{Start: Point{X: 1, Y: 1}}
	for y := 1; y <= 3; y++ {
		for x := 1; x <= 3; x++ {
			maze.Tiles[y][x] = Floor
		}
	}
	path, found := maze.FindPath(Point{X: 1, Y: 1}, Point{X: 3, Y: 3})
	if !found || len(path) != 4 || path[len(path)-1] != (Point{X: 3, Y: 3}) {
		t.Fatalf("path = %v, found = %t; want 4-step route to goal", path, found)
	}
	for index := 1; index < len(path); index++ {
		previous := path[index-1]
		current := path[index]
		distance := abs(previous.X-current.X) + abs(previous.Y-current.Y)
		if distance != 1 {
			t.Fatalf("path has non-adjacent steps: %v", path)
		}
	}
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
