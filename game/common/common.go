package common

const (
	GreeedSize      = 100
	ScreenWidth     = 1100.0
	ScreenHeight    = 900.0
	ScreenGreedSize = 50
	PlayerSize      = 50
	PlayerHeadSize  = 50

	SwordAttachmentWidth     = 195
	SwordAttachmentHeight    = 65
	SwordAttachmentSmoothing = 0.6
)

var RoomWidth = 1500.0
var RoomHeight = 1500.0

var MaxPlayerSpeedMoving = 1000

var HeadAcceleration = 1000.0
var HeadNormalDeceleration = 800.0
var HeadDeceleration = HeadNormalDeceleration

var Acceleration = 500.0
var Deceleration = 300.0

var PlayerDamage int
var PlayerHelth = MaxPlayerHelth
var MaxPlayerHelth = 100.0
var MaxSpeed = 300.0
var HeadRotationSpeed = 450.0
var Score int = 1

var SwordExist bool

type PointPlayer struct {
	Px, Py float64
}

type PointSpeed struct {
	Vx, Vy float64
}
