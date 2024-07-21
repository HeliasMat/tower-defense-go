package entities

import (
	"math"
	"testing"
	"tower-defense/internal/entities"
)

func TestNewEnemy(t *testing.T) {
	path := []entities.BaseEntity{
		{X: 0, Y: 0},
		{X: 1, Y: 1},
	}
	e := entities.NewEnemy(100, 10, 5, 1.0, path)

	if e.Health != 100 || e.MaxHealth != 100 {
		t.Errorf("Expected Health and MaxHealth to be 100, got %d and %d", e.Health, e.MaxHealth)
	}
	if e.Reward != 10 {
		t.Errorf("Expected Reward to be 10, got %d", e.Reward)
	}
	if e.Damage != 5 {
		t.Errorf("Expected Damage to be 5, got %d", e.Damage)
	}
	if e.Speed != 1.0 {
		t.Errorf("Expected Speed to be 1.0, got %f", e.Speed)
	}
	if e.PathIndex != 0 {
		t.Errorf("Expected PathIndex to be 0, got %d", e.PathIndex)
	}
	if len(e.Path) != 2 {
		t.Errorf("Expected Path length to be 2, got %d", len(e.Path))
	}
	if e.X != 0 || e.Y != 0 {
		t.Errorf("Expected initial position to be (0,0), got (%f,%f)", e.X, e.Y)
	}
}

func TestTakeDamage(t *testing.T) {
	e := entities.NewEnemy(100, 10, 5, 1.0, []entities.BaseEntity{{X: 0, Y: 0}})

	tests := []struct {
		damage   int
		expected int
		dead     bool
	}{
		{50, 50, false},
		{30, 20, false},
		{30, 0, true},
	}

	for _, tt := range tests {
		dead := e.TakeDamage(tt.damage)
		if e.Health != tt.expected {
			t.Errorf("Expected Health to be %d after %d damage, got %d", tt.expected, tt.damage, e.Health)
		}
		if dead != tt.dead {
			t.Errorf("Expected dead to be %v after %d damage, got %v", tt.dead, tt.damage, dead)
		}
	}

	// Test that health doesn't go below 0
	e.Health = 10
	e.TakeDamage(20)
	if e.Health != 0 {
		t.Errorf("Expected Health to be 0 after excessive damage, got %d", e.Health)
	}
}

func TestMove(t *testing.T) {
	path := []entities.BaseEntity{
		{X: 0, Y: 0},
		{X: 3, Y: 4},
		{X: 6, Y: 8},
	}
	e := entities.NewEnemy(100, 10, 5, 2.0, path)

	// Test partial movement
	e.Move()
	if math.Abs(e.X-1.8) > 0.001 || math.Abs(e.Y-2.4) > 0.001 {
		t.Errorf("Expected position after partial move to be (1.8,2.4), got (%f,%f)", e.X, e.Y)
	}
	if e.PathIndex != 0 {
		t.Errorf("Expected PathIndex to still be 0 after partial move, got %d", e.PathIndex)
	}

	// Test reaching next point
	e.Move()
	if e.X != 3 || e.Y != 4 {
		t.Errorf("Expected position after full move to be (3,4), got (%f,%f)", e.X, e.Y)
	}
	if e.PathIndex != 1 {
		t.Errorf("Expected PathIndex to be 1 after reaching next point, got %d", e.PathIndex)
	}

	// Test staying at end of path
	e.PathIndex = len(path) - 1
	initialX, initialY := e.X, e.Y
	e.Move()
	if e.X != initialX || e.Y != initialY {
		t.Errorf("Expected no movement at end of path, got movement from (%f,%f) to (%f,%f)", initialX, initialY, e.X, e.Y)
	}
}

func TestHasReachedEnd(t *testing.T) {
	path := []entities.BaseEntity{
		{X: 0, Y: 0},
		{X: 1, Y: 1},
	}
	e := entities.NewEnemy(100, 10, 5, 1.0, path)

	if e.HasReachedEnd() {
		t.Error("Expected HasReachedEnd to be false at start of path")
	}

	e.PathIndex = len(path) - 1
	if !e.HasReachedEnd() {
		t.Error("Expected HasReachedEnd to be true at end of path")
	}
}

func TestGetReward(t *testing.T) {
	e := entities.NewEnemy(100, 10, 5, 1.0, []entities.BaseEntity{{X: 0, Y: 0}})
	if e.GetReward() != 10 {
		t.Errorf("Expected GetReward to return 10, got %d", e.GetReward())
	}
}

func TestGetDamage(t *testing.T) {
	e := entities.NewEnemy(100, 10, 5, 1.0, []entities.BaseEntity{{X: 0, Y: 0}})
	if e.GetDamage() != 5 {
		t.Errorf("Expected GetDamage to return 5, got %d", e.GetDamage())
	}
}
