package entities

import (
	"testing"
	"time"
	"tower-defense/internal/entities"
)

func TestNewTowers(t *testing.T) {
	tests := []struct {
		name     string
		newFunc  func(float64, float64) *entities.Tower
		expected entities.Tower
	}{
		{
			name:    "BasicTower",
			newFunc: entities.NewBasicTower,
			expected: entities.Tower{
				BaseEntity: entities.BaseEntity{X: 10, Y: 20},
				Range:      100,
				Damage:     10,
				FireRate:   time.Second,
				Level:      1,
				Cost:       50,
				Type:       "Basic",
			},
		},
		{
			name:    "SniperTower",
			newFunc: entities.NewSniperTower,
			expected: entities.Tower{
				BaseEntity: entities.BaseEntity{X: 10, Y: 20},
				Range:      200,
				Damage:     30,
				FireRate:   time.Second * 2,
				Level:      1,
				Cost:       100,
				Type:       "Sniper",
			},
		},
		{
			name:    "AOETower",
			newFunc: entities.NewAOETower,
			expected: entities.Tower{
				BaseEntity: entities.BaseEntity{X: 10, Y: 20},
				Range:      80,
				Damage:     15,
				FireRate:   time.Second * 2,
				Level:      1,
				Cost:       150,
				Type:       "AOE",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tower := tt.newFunc(10, 20)
			if tower.X != tt.expected.X || tower.Y != tt.expected.Y {
				t.Errorf("Expected position (%f,%f), got (%f,%f)", tt.expected.X, tt.expected.Y, tower.X, tower.Y)
			}
			if tower.Range != tt.expected.Range {
				t.Errorf("Expected Range %f, got %f", tt.expected.Range, tower.Range)
			}
			if tower.Damage != tt.expected.Damage {
				t.Errorf("Expected Damage %d, got %d", tt.expected.Damage, tower.Damage)
			}
			if tower.FireRate != tt.expected.FireRate {
				t.Errorf("Expected FireRate %v, got %v", tt.expected.FireRate, tower.FireRate)
			}
			if tower.Level != tt.expected.Level {
				t.Errorf("Expected Level %d, got %d", tt.expected.Level, tower.Level)
			}
			if tower.Cost != tt.expected.Cost {
				t.Errorf("Expected Cost %d, got %d", tt.expected.Cost, tower.Cost)
			}
			if tower.Type != tt.expected.Type {
				t.Errorf("Expected Type %s, got %s", tt.expected.Type, tower.Type)
			}
		})
	}
}

func TestCanFire(t *testing.T) {
	tower := entities.NewBasicTower(0, 0)

	if !tower.CanFire() {
		t.Error("Expected CanFire to be true initially")
	}

	tower.Fire()
	if tower.CanFire() {
		t.Error("Expected CanFire to be false immediately after firing")
	}

	time.Sleep(tower.FireRate + time.Millisecond)
	if !tower.CanFire() {
		t.Error("Expected CanFire to be true after FireRate duration")
	}
}

func TestUpgrade(t *testing.T) {
	tower := entities.NewBasicTower(0, 0)
	initialDamage := tower.Damage
	initialRange := tower.Range
	initialFireRate := tower.FireRate

	err := tower.Upgrade()
	if err != nil {
		t.Errorf("Unexpected error on first upgrade: %v", err)
	}
	if tower.Level != 2 {
		t.Errorf("Expected Level 2 after upgrade, got %d", tower.Level)
	}
	if tower.Damage != initialDamage+5 {
		t.Errorf("Expected Damage to increase by 5, got %d", tower.Damage)
	}
	if tower.Range != initialRange+20 {
		t.Errorf("Expected Range to increase by 20, got %f", tower.Range)
	}
	if tower.FireRate != time.Duration(float64(initialFireRate)*0.9) {
		t.Errorf("Expected FireRate to decrease to 90%%, got %v", tower.FireRate)
	}

	tower.Upgrade()
	err = tower.Upgrade()
	if err == nil {
		t.Error("Expected error on upgrading past maximum level")
	}
}

func TestGetUpgradeCost(t *testing.T) {
	tower := entities.NewBasicTower(0, 0)
	if tower.GetUpgradeCost() != tower.Cost {
		t.Errorf("Expected upgrade cost %d, got %d", tower.Cost, tower.GetUpgradeCost())
	}

	tower.Upgrade()
	if tower.GetUpgradeCost() != tower.Cost*2 {
		t.Errorf("Expected upgrade cost %d, got %d", tower.Cost*2, tower.GetUpgradeCost())
	}
}

func TestGetSellValue(t *testing.T) {
	tower := entities.NewBasicTower(0, 0)
	if tower.GetSellValue() != tower.Cost/2 {
		t.Errorf("Expected sell value %d, got %d", tower.Cost/2, tower.GetSellValue())
	}

	tower.Upgrade()
	if tower.GetSellValue() != tower.Cost {
		t.Errorf("Expected sell value %d, got %d", tower.Cost, tower.GetSellValue())
	}
}

func TestUpdate(t *testing.T) {
	tower := entities.NewBasicTower(0, 0)
	enemy1 := entities.NewEnemy(100, 10, 5, 1.0, []entities.BaseEntity{{X: 50, Y: 0}})
	enemy2 := entities.NewEnemy(100, 10, 5, 1.0, []entities.BaseEntity{{X: 200, Y: 0}})
	enemies := []*entities.Enemy{enemy1, enemy2}

	tower.Update(enemies)
	if enemy1.Health != 90 {
		t.Errorf("Expected enemy1 Health to be 90, got %d", enemy1.Health)
	}
	if enemy2.Health != 100 {
		t.Errorf("Expected enemy2 Health to be 100, got %d", enemy2.Health)
	}

	tower.Update(enemies) // Should not fire due to fire rate
	if enemy1.Health != 90 {
		t.Errorf("Expected enemy1 Health to still be 90, got %d", enemy1.Health)
	}
}

func TestIsInRange(t *testing.T) {
	tower := entities.NewBasicTower(0, 0)
	enemy1 := entities.NewEnemy(100, 10, 5, 1.0, []entities.BaseEntity{{X: 50, Y: 0}})
	enemy2 := entities.NewEnemy(100, 10, 5, 1.0, []entities.BaseEntity{{X: 150, Y: 0}})

	if !tower.IsInRange(enemy1) {
		t.Error("Expected enemy1 to be in range")
	}
	if tower.IsInRange(enemy2) {
		t.Error("Expected enemy2 to be out of range")
	}
}

func TestDealAOEDamage(t *testing.T) {
	tower := entities.NewAOETower(0, 0)
	target := entities.NewEnemy(100, 10, 5, 1.0, []entities.BaseEntity{{X: 50, Y: 0}})
	enemy1 := entities.NewEnemy(100, 10, 5, 1.0, []entities.BaseEntity{{X: 60, Y: 0}})
	enemy2 := entities.NewEnemy(100, 10, 5, 1.0, []entities.BaseEntity{{X: 100, Y: 0}})
	enemies := []*entities.Enemy{target, enemy1, enemy2}

	tower.DealAOEDamage(enemies, target)

	if target.Health != 100 {
		t.Errorf("Expected target Health to be 100, got %d", target.Health)
	}
	if enemy1.Health != 100-tower.Damage/2 {
		t.Errorf("Expected enemy1 Health to be %d, got %d", 100-tower.Damage/2, enemy1.Health)
	}
	if enemy2.Health != 100 {
		t.Errorf("Expected enemy2 Health to be 100, got %d", enemy2.Health)
	}
}
