package core

import (
	"testing"
	"tower-defense/internal/core"
	"tower-defense/internal/entities"
)

func TestNewGameState(t *testing.T) {
	gs := core.NewGameState()

	if len(gs.GetTowers()) != 0 {
		t.Errorf("Expected 0 towers, got %d", len(gs.GetTowers()))
	}
	if len(gs.GetEnemies()) != 0 {
		t.Errorf("Expected 0 enemies, got %d", len(gs.GetEnemies()))
	}
	if gs.GetLives() != 100 {
		t.Errorf("Expected 100 lives, got %d", gs.GetLives())
	}
	if gs.GetMoney() != 1000 {
		t.Errorf("Expected 1000 money, got %d", gs.GetMoney())
	}
	if gs.GetWave() != 0 {
		t.Errorf("Expected wave 0, got %d", gs.GetWave())
	}
	if len(gs.GetTowerCosts()) != 3 {
		t.Errorf("Expected 3 tower types, got %d", len(gs.GetTowerCosts()))
	}
	if gs.IsPaused() {
		t.Error("Expected game to start unpaused")
	}
	if len(gs.GetEnemyPath()) != 8 {
		t.Errorf("Expected 8 path points, got %d", len(gs.GetEnemyPath()))
	}
}

func TestAddTower(t *testing.T) {
	gs := core.NewGameState()

	tests := []struct {
		name      string
		towerType core.TowerType
		x, y      float64
		expectErr bool
	}{
		{"Add Basic Tower", core.BasicTower, 100, 100, false},
		{"Add Sniper Tower", core.SniperTower, 200, 200, false},
		{"Add AOE Tower", core.AOETower, 300, 300, false},
		{"Invalid Tower Type", core.TowerType(999), 400, 400, true},
		{"Not Enough Money", core.BasicTower, 500, 500, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialMoney := gs.GetMoney()
			initialTowerCount := len(gs.GetTowers())
			err := gs.AddTower(tt.towerType, tt.x, tt.y)

			if tt.expectErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				if len(gs.GetTowers()) != initialTowerCount {
					t.Error("Tower count should not change on error")
				}
				if gs.GetMoney() != initialMoney {
					t.Error("Money should not change on error")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(gs.GetTowers()) != initialTowerCount+1 {
					t.Error("Tower count should increase by 1")
				}
				if gs.GetMoney() >= initialMoney {
					t.Error("Money should decrease")
				}
			}
		})
	}
}

func TestAddAndRemoveEnemy(t *testing.T) {
	gs := core.NewGameState()
	enemy := entities.NewEnemy(100, 10, 1, 1.0, gs.GetEnemyPath())

	gs.AddEnemy(enemy)
	if len(gs.GetEnemies()) != 1 {
		t.Errorf("Expected 1 enemy, got %d", len(gs.GetEnemies()))
	}

	gs.RemoveEnemy(0)
	if len(gs.GetEnemies()) != 0 {
		t.Errorf("Expected 0 enemies after removal, got %d", len(gs.GetEnemies()))
	}

	gs.RemoveEnemy(0) // Should not panic
}

func TestDamageEnemy(t *testing.T) {
	gs := core.NewGameState()
	enemy := entities.NewEnemy(100, 10, 1, 1.0, gs.GetEnemyPath())
	gs.AddEnemy(enemy)

	initialMoney := gs.GetMoney()
	isDead := gs.DamageEnemy(0, 50)
	if isDead {
		t.Error("Enemy should not be dead after 50 damage")
	}
	if gs.GetMoney() != initialMoney {
		t.Error("Money should not change when enemy is not killed")
	}

	isDead = gs.DamageEnemy(0, 60)
	if !isDead {
		t.Error("Enemy should be dead after 110 total damage")
	}
	if gs.GetMoney() != initialMoney+enemy.GetReward() {
		t.Error("Money should increase by enemy reward when killed")
	}
	if len(gs.GetEnemies()) != 0 {
		t.Error("Enemy should be removed after being killed")
	}
}

func TestLoseLife(t *testing.T) {
	gs := core.NewGameState()
	initialLives := gs.GetLives()

	gs.LoseLife(10)
	if gs.GetLives() != initialLives-10 {
		t.Errorf("Expected %d lives, got %d", initialLives-10, gs.GetLives())
	}

	gs.LoseLife(1000)
	if gs.GetLives() != 0 {
		t.Error("Lives should not go below 0")
	}
}

func TestIsGameOver(t *testing.T) {
	gs := core.NewGameState()
	if gs.IsGameOver() {
		t.Error("Game should not be over at start")
	}

	gs.LoseLife(100)
	if !gs.IsGameOver() {
		t.Error("Game should be over when lives reach 0")
	}
}

func TestNextWave(t *testing.T) {
	gs := core.NewGameState()
	initialWave := gs.GetWave()
	initialEnemyCount := len(gs.GetEnemies())

	gs.NextWave()
	if gs.GetWave() != initialWave+1 {
		t.Errorf("Expected wave %d, got %d", initialWave+1, gs.GetWave())
	}
	if len(gs.GetEnemies()) <= initialEnemyCount {
		t.Error("Enemy count should increase after next wave")
	}
}

func TestUpgradeTower(t *testing.T) {
	gs := core.NewGameState()
	gs.AddTower(core.BasicTower, 100, 100)

	initialMoney := gs.GetMoney()
	err := gs.UpgradeTower(0)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if gs.GetMoney() >= initialMoney {
		t.Error("Money should decrease after upgrade")
	}
	if gs.GetTowers()[0].Level != 2 {
		t.Errorf("Expected tower level 2, got %d", gs.GetTowers()[0].Level)
	}

	gs.SetMoney(0)
	err = gs.UpgradeTower(0)
	if err == nil {
		t.Error("Expected error when upgrading without enough money")
	}
}

func TestSellTower(t *testing.T) {
	gs := core.NewGameState()
	gs.AddTower(core.BasicTower, 100, 100)

	initialMoney := gs.GetMoney()
	initialTowerCount := len(gs.GetTowers())

	err := gs.SellTower(0)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if gs.GetMoney() <= initialMoney {
		t.Error("Money should increase after selling tower")
	}
	if len(gs.GetTowers()) != initialTowerCount-1 {
		t.Error("Tower count should decrease after selling")
	}

	err = gs.SellTower(0)
	if err == nil {
		t.Error("Expected error when selling non-existent tower")
	}
}

func TestTogglePause(t *testing.T) {
	gs := core.NewGameState()
	if gs.IsPaused() {
		t.Error("Game should start unpaused")
	}

	gs.TogglePause()
	if !gs.IsPaused() {
		t.Error("Game should be paused after toggle")
	}

	gs.TogglePause()
	if gs.IsPaused() {
		t.Error("Game should be unpaused after second toggle")
	}
}

func TestUpdate(t *testing.T) {
	gs := core.NewGameState()
	enemy := entities.NewEnemy(100, 10, 1, 1000.0, gs.GetEnemyPath())
	gs.AddEnemy(enemy)
	gs.AddTower(core.BasicTower, gs.GetEnemyPath()[1].X, gs.GetEnemyPath()[1].Y)

	initialLives := gs.GetLives()
	initialMoney := gs.GetMoney()

	gs.Update()
	if len(gs.GetEnemies()) != 0 {
		t.Error("Enemy should reach end and be removed")
	}
	if gs.GetLives() != initialLives-enemy.GetDamage() {
		t.Error("Lives should decrease when enemy reaches end")
	}
	if gs.GetMoney() != initialMoney {
		t.Error("Money should not change from update alone")
	}

	gs.SetPaused(true)
	gs.AddEnemy(enemy)
	gs.Update()
	if len(gs.GetEnemies()) != 1 {
		t.Error("No updates should occur while paused")
	}
}

func TestGetterMethods(t *testing.T) {
	gs := core.NewGameState()
	gs.AddTower(core.BasicTower, 100, 100)
	gs.AddEnemy(entities.NewEnemy(100, 10, 1, 1.0, gs.GetEnemyPath()))

	if len(gs.GetTowers()) != 1 {
		t.Error("GetTowers should return 1 tower")
	}
	if len(gs.GetEnemies()) != 1 {
		t.Error("GetEnemies should return 1 enemy")
	}
	if gs.GetLives() != 100 {
		t.Error("GetLives should return 100")
	}
	if gs.GetMoney() != 950 {
		t.Error("GetMoney should return 950")
	}
	if gs.GetWave() != 0 {
		t.Error("GetWave should return 0")
	}
	if len(gs.GetTowerCosts()) != 3 {
		t.Error("GetTowerCosts should return 3 tower types")
	}
	if gs.IsPaused() {
		t.Error("IsPaused should return false initially")
	}
	if len(gs.GetEnemyPath()) != 8 {
		t.Error("GetEnemyPath should return 8 path points")
	}
}
