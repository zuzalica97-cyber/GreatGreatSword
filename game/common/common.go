package common

const (
	ScreenWidth  = 800
	ScreenHeight = 600
	PlayerSize   = 25

	MaxSpeed     = 500
	Acceleration = 400
	Deceleration = 600
)

type PointPlayer struct {
	Px, Py float64
}

type PointSpeed struct {
	Vx, Vy float64
}
