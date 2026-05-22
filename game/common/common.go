package common

const (
	ScreenWidth    = 800
	ScreenHeight   = 600
	PlayerSize     = 25
	PlayerHeadSize = 25

	MaxSpeed     = 300.0
	Acceleration = 400.0
	Deceleration = 600.0

	HeadRotationSpeed = 360.0
	HeadAcceleration  = 800.0
	HeadDeceleration  = 300.0

	SwordSizeLower           = 60
	SwordSizeUpper           = 30
	SwordAttachmentSmoothing = 0.85
)

type PointPlayer struct {
	Px, Py float64
}

type PointSpeed struct {
	Vx, Vy float64
}
