package common

const (
	ScreenWidth    = 1100
	ScreenHeight   = 900
	PlayerSize     = 40
	PlayerHeadSize = 40

	Acceleration = 500.0
	Deceleration = 300.0

	HeadAcceleration = 500.0
	HeadDeceleration = 200.0

	SwordAttachmentWidth     = 130
	SwordAttachmentHeight    = 50
	SwordAttachmentSmoothing = 0.6
)

var MaxSpeed = 300.0
var HeadRotationSpeed = 450.0

type PointPlayer struct {
	Px, Py float64
}

type PointSpeed struct {
	Vx, Vy float64
}
