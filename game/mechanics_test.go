package game

import (
	"math/rand"
	"testing"
)

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
	if current.State != StateDefeat || current.Player.Energy != 0 {
		t.Fatalf("zero-energy move state = %v, energy = %d", current.State, current.Player.Energy)
	}
}

func TestArmedCombatBreaksSwordAndRemovesZog(t *testing.T) {
	current := NewGameWithSeed(12)
	current.Player.HasSword = true
	current.Maze.Tiles[1][2] = Monster
	current.Move(1, 0)
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
	if current.State != StateDefeat || current.StatusMessage != "The Zog overpowered you!" {
		t.Fatalf("unarmed combat state = %v, message = %q", current.State, current.StatusMessage)
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
	if current.State != StatePlaying || current.ZogsSlain != 1 || current.Player.Energy != InitialEnergy-41 {
		t.Fatalf("unarmed win state = %v, kills = %d, energy = %d", current.State, current.ZogsSlain, current.Player.Energy)
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

	current.Interact()
	if current.Player.Energy != InitialEnergy-15 || current.GuidanceTicks != 15 || len(current.GuidancePath) == 0 {
		t.Fatalf("interaction energy=%d ticks=%d path length=%d", current.Player.Energy, current.GuidanceTicks, len(current.GuidancePath))
	}
	if current.GuidancePath[len(current.GuidancePath)-1] != current.Maze.Treasure {
		t.Fatal("guidance path does not end at the treasure")
	}

	current.Player.HasTreasure = true
	current.Interact()
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
