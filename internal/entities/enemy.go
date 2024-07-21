package entities

import "math"

type Enemy struct {
	BaseEntity
	Health    int
	MaxHealth int
	Speed     float64
	Reward    int
	Damage    int
	PathIndex int
	Path      []BaseEntity
}

func NewEnemy(health, reward, damage int, speed float64, path []BaseEntity) *Enemy {
	return &Enemy{
		BaseEntity: BaseEntity{X: path[0].X, Y: path[0].Y},
		Health:     health,
		MaxHealth:  health,
		Speed:      speed,
		Reward:     reward,
		Damage:     damage,
		PathIndex:  0,
		Path:       path,
	}
}

func (e *Enemy) TakeDamage(damage int) bool {
	e.Health -= damage
	if e.Health < 0 {
		e.Health = 0
	}
	return e.Health <= 0
}

func (e *Enemy) Move() {
	if e.PathIndex >= len(e.Path)-1 {
		return
	}

	target := e.Path[e.PathIndex+1]
	dx := target.X - e.X
	dy := target.Y - e.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance <= e.Speed {
		e.X = target.X
		e.Y = target.Y
		e.PathIndex++
	} else {
		e.X += (dx / distance) * e.Speed
		e.Y += (dy / distance) * e.Speed
	}
}

func (e *Enemy) HasReachedEnd() bool {
	return e.PathIndex >= len(e.Path)-1
}

func (e *Enemy) GetReward() int {
	return e.Reward
}

func (e *Enemy) GetDamage() int {
	return e.Damage
}
