package core

import (
	"errors"
	"sync"
	"tower-defense/internal/entities"
)

type TowerType int

const (
	BasicTower TowerType = iota
	SniperTower
	AOETower
)

type GameState struct {
	mu         sync.RWMutex
	towers     []*entities.Tower
	enemies    []*entities.Enemy
	lives      int
	money      int
	wave       int
	towerCosts map[TowerType]int
	paused     bool
	enemyPath  []entities.BaseEntity
}

func NewGameState() *GameState {
	return &GameState{
		towers:  make([]*entities.Tower, 0, 100), // Pre-allocate space for 100 towers
		enemies: make([]*entities.Enemy, 0, 200), // Pre-allocate space for 200 enemies
		lives:   100,
		money:   1000,
		wave:    0,
		towerCosts: map[TowerType]int{
			BasicTower:  50,
			SniperTower: 100,
			AOETower:    150,
		},
		paused: false,
		enemyPath: []entities.BaseEntity{
			{X: 0, Y: 300},
			{X: 200, Y: 300},
			{X: 200, Y: 100},
			{X: 400, Y: 100},
			{X: 400, Y: 500},
			{X: 600, Y: 500},
			{X: 600, Y: 300},
			{X: 800, Y: 300},
		},
	}
}

func (gs *GameState) AddTower(towerType TowerType, x, y float64) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	cost, exists := gs.towerCosts[towerType]
	if !exists {
		return errors.New("invalid tower type")
	}
	if gs.money < cost {
		return errors.New("not enough money to add tower")
	}

	var tower *entities.Tower
	switch towerType {
	case BasicTower:
		tower = entities.NewBasicTower(x, y)
	case SniperTower:
		tower = entities.NewSniperTower(x, y)
	case AOETower:
		tower = entities.NewAOETower(x, y)
	default:
		return errors.New("unknown tower type")
	}

	gs.towers = append(gs.towers, tower)
	gs.money -= cost
	return nil
}

func (gs *GameState) AddEnemy(enemy *entities.Enemy) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.enemies = append(gs.enemies, enemy)
}

func (gs *GameState) RemoveEnemy(index int) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	if index < 0 || index >= len(gs.enemies) {
		return
	}
	gs.enemies[index] = gs.enemies[len(gs.enemies)-1]
	gs.enemies = gs.enemies[:len(gs.enemies)-1]
}

func (gs *GameState) DamageEnemy(index int, damage int) bool {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	if index < 0 || index >= len(gs.enemies) {
		return false
	}
	enemy := gs.enemies[index]
	isDead := enemy.TakeDamage(damage)
	if isDead {
		gs.money += enemy.GetReward()
		gs.enemies[index] = gs.enemies[len(gs.enemies)-1]
		gs.enemies = gs.enemies[:len(gs.enemies)-1]
	}
	return isDead
}

func (gs *GameState) LoseLife(amount int) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.lives -= amount
	if gs.lives < 0 {
		gs.lives = 0
	}
}

func (gs *GameState) IsGameOver() bool {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.lives <= 0
}

func (gs *GameState) NextWave() {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.wave++
	gs.spawnEnemiesForWave()
}

func (gs *GameState) spawnEnemiesForWave() {
	numEnemies := gs.wave * 2 // Example: 2 enemies per wave
	for i := 0; i < numEnemies; i++ {
		health := 50 + gs.wave*10
		speed := 1.0 + float64(gs.wave)/10.0
		reward := 10 + gs.wave
		damage := 1 + gs.wave/5
		enemy := entities.NewEnemy(health, reward, damage, speed, gs.enemyPath)
		gs.enemies = append(gs.enemies, enemy)
	}
}

func (gs *GameState) UpgradeTower(index int) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	if index < 0 || index >= len(gs.towers) {
		return errors.New("invalid tower index")
	}
	tower := gs.towers[index]
	upgradeCost := tower.GetUpgradeCost()
	if gs.money < upgradeCost {
		return errors.New("not enough money to upgrade tower")
	}
	if err := tower.Upgrade(); err != nil {
		return err
	}
	gs.money -= upgradeCost
	return nil
}

func (gs *GameState) SellTower(index int) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	if index < 0 || index >= len(gs.towers) {
		return errors.New("invalid tower index")
	}
	tower := gs.towers[index]
	sellValue := tower.GetSellValue()
	gs.money += sellValue
	gs.towers[index] = gs.towers[len(gs.towers)-1]
	gs.towers = gs.towers[:len(gs.towers)-1]
	return nil
}

func (gs *GameState) TogglePause() {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.paused = !gs.paused
}

func (gs *GameState) Update() {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	if gs.paused {
		return
	}

	for _, tower := range gs.towers {
		tower.Update(gs.enemies)
	}

	for i := 0; i < len(gs.enemies); i++ {
		enemy := gs.enemies[i]
		if enemy.HasReachedEnd() {
			gs.lives -= enemy.GetDamage()
			if gs.lives < 0 {
				gs.lives = 0
			}
			gs.enemies[i] = gs.enemies[len(gs.enemies)-1]
			gs.enemies = gs.enemies[:len(gs.enemies)-1]
			i--
		} else {
			enemy.Move()
		}
	}

	if len(gs.enemies) == 0 {
		gs.wave++
		gs.spawnEnemiesForWave()
	}
}

// Getter methods for private fields
func (gs *GameState) GetTowers() []*entities.Tower {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.towers
}

func (gs *GameState) GetEnemies() []*entities.Enemy {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.enemies
}

func (gs *GameState) GetLives() int {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.lives
}

func (gs *GameState) GetMoney() int {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.money
}

func (gs *GameState) GetWave() int {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.wave
}

func (gs *GameState) GetTowerCosts() map[TowerType]int {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.towerCosts
}

func (gs *GameState) IsPaused() bool {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.paused
}

func (gs *GameState) GetEnemyPath() []entities.BaseEntity {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.enemyPath
}

func (gs *GameState) SetTowers(towers []*entities.Tower) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.towers = towers
}

func (gs *GameState) SetEnemies(enemies []*entities.Enemy) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.enemies = enemies
}

func (gs *GameState) SetLives(lives int) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.lives = lives
}

func (gs *GameState) SetMoney(money int) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.money = money
}

func (gs *GameState) SetWave(wave int) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.wave = wave
}

func (gs *GameState) SetTowerCosts(towerCosts map[TowerType]int) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.towerCosts = towerCosts
}

func (gs *GameState) SetPaused(paused bool) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.paused = paused
}

func (gs *GameState) SetEnemyPath(enemyPath []entities.BaseEntity) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.enemyPath = enemyPath
}
