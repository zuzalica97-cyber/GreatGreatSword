package common

const (
	ScreenWidth    = 1100
	ScreenHeight   = 900
	PlayerSize     = 50
	PlayerHeadSize = 50

	Acceleration = 500.0
	Deceleration = 300.0

	HeadAcceleration = 1000.0
	HeadDeceleration = 800.0

	SwordAttachmentWidth     = 165
	SwordAttachmentHeight    = 65
	SwordAttachmentSmoothing = 0.6

	PatheticSize      = 60
	PatheticBaseSpeed = 50
	PatheticMaxSpeed  = 300
	PatheticDistanse  = 450
	PatheticDamage    = 15

	HaterSize           = 50
	HaterBaseSpeed      = 50
	HaterMaxSpeed       = 300
	HaterDistanse       = 300
	HaitersBolletDamage = 10
)

var PlayerHelth = MaxPlayerHelth
var MaxPlayerHelth = 100
var MaxSpeed = 300.0
var HeadRotationSpeed = 450.0
var Score int = 0

type PointPlayer struct {
	Px, Py float64
}

type PointSpeed struct {
	Vx, Vy float64
}
