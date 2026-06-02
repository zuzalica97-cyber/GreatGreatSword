package common

const (
	ScreenWidth    = 1100
	ScreenHeight   = 900
	PlayerSize     = 50
	PlayerHeadSize = 50

	SwordAttachmentWidth     = 165
	SwordAttachmentHeight    = 65
	SwordAttachmentSmoothing = 0.6

	PatheticSize     = 60
	PatheticMaxSpeed = 300
	PatheticDistanse = 450

	HaterSize           = 50
	HaterBaseSpeed      = 50
	HaterMaxSpeed       = 300
	HaterDistanse       = 300
	HaitersBolletDamage = 10
)

var MaxPlayerSpeedMoving = 1000

var HeadAcceleration = 1000.0
var HeadNormalDeceleration = 800.0
var HeadDeceleration = HeadNormalDeceleration

var PatheticNormalSpeed = 50.0
var PatheticBaseSpeed = PatheticNormalSpeed

var PathiticNormalDamage = 5
var PatheticDamage = PathiticNormalDamage

var NormalAcceleration = 500.0
var NormalDeceleration = 300.0

var Acceleration = NormalAcceleration
var Deceleration = NormalAcceleration

var PlayerHelth = MaxPlayerHelth
var MaxPlayerHelth = 100
var MaxSpeed = 300.0
var HeadRotationSpeed = 450.0
var Score int = 1

var MaxValwe int = 15
var Valwe int = 3

type PointPlayer struct {
	Px, Py float64
}

type PointSpeed struct {
	Vx, Vy float64
}
