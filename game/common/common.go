package common

const (
	ScreenWidth    = 1100
	ScreenHeight   = 900
	PlayerSize     = 50
	PlayerHeadSize = 50

	Acceleration = 500.0
	Deceleration = 300.0

	HeadAcceleration = 500.0
	HeadDeceleration = 800.0

	SwordAttachmentWidth     = 165
	SwordAttachmentHeight    = 65
	SwordAttachmentSmoothing = 0.6

	PatheticSize      = 60
	PatheticBaseSpeed = 50
	PatheticMaxSpeed  = 300
	PatheticDistanse  = 450
)

var MaxSpeed = 300.0
var HeadRotationSpeed = 450.0
var Score int = 0

type PointPlayer struct {
	Px, Py float64
}

type PointSpeed struct {
	Vx, Vy float64
}
