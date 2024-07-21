package entities

type Entity interface {
	GetPosition() (float64, float64)
	SetPosition(x, y float64)
}

type BaseEntity struct {
	X, Y float64
}

func (e *BaseEntity) GetPosition() (float64, float64) {
	return e.X, e.Y
}

func (e *BaseEntity) SetPosition(x, y float64) {
	e.X, e.Y = x, y
}
