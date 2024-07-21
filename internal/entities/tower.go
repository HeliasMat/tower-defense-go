package entities

import (
	"errors"
	"time"
)

type Tower struct {
	BaseEntity
	Range     float64
	Damage    int
	FireRate  time.Duration
	LastFired time.Time
	Level     int
	Cost      int
	Type      string
}

func NewBasicTower(x, y float64) *Tower {
	return &Tower{
		BaseEntity: BaseEntity{X: x, Y: y},
		Range:      100,
		Damage:     10,
		FireRate:   time.Second,
		Level:      1,
		Cost:       50,
		Type:       "Basic",
	}
}

func NewSniperTower(x, y float64) *Tower {
	return &Tower{
		BaseEntity: BaseEntity{X: x, Y: y},
		Range:      200,
		Damage:     30,
		FireRate:   time.Second * 2,
		Level:      1,
		Cost:       100,
		Type:       "Sniper",
	}
}

func NewAOETower(x, y float64) *Tower {
	return &Tower{
		BaseEntity: BaseEntity{X: x, Y: y},
		Range:      80,
		Damage:     15,
		FireRate:   time.Second * 2,
		Level:      1,
		Cost:       150,
		Type:       "AOE",
	}
}

func (t *Tower) CanFire() bool {
	return time.Since(t.LastFired) >= t.FireRate
}

func (t *Tower) Fire() {
	t.LastFired = time.Now()
}

func (t *Tower) Upgrade() error {
	if t.Level >= 3 {
		return errors.New("tower is already at maximum level")
	}
	t.Level++
	t.Damage += 5
	t.Range += 20
	t.FireRate = time.Duration(float64(t.FireRate) * 0.9)
	return nil
}

func (t *Tower) GetUpgradeCost() int {
	return t.Cost * t.Level
}

func (t *Tower) GetSellValue() int {
	return t.Cost * t.Level / 2
}

func (t *Tower) Update(enemies []*Enemy) {
	if !t.CanFire() {
		return
	}

	for _, enemy := range enemies {
		if t.IsInRange(enemy) {
			t.Fire()
			enemy.TakeDamage(t.Damage)
			if t.Type == "AOE" {
				t.DealAOEDamage(enemies, enemy)
			}
			break
		}
	}
}

func (t *Tower) IsInRange(e *Enemy) bool {
	dx := t.X - e.X
	dy := t.Y - e.Y
	distanceSquared := dx*dx + dy*dy
	return distanceSquared <= t.Range*t.Range
}

func (t *Tower) DealAOEDamage(enemies []*Enemy, target *Enemy) {
	for _, enemy := range enemies {
		if enemy != target && t.IsInRange(enemy) {
			enemy.TakeDamage(t.Damage / 2) // AOE damage is half of the main target
		}
	}
}
