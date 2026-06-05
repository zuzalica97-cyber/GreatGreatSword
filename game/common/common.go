package common

const (
	ScreenWidth     = 1100
	ScreenHeight    = 900
	ScreenGreedSize = 50
	PlayerSize      = 50
	PlayerHeadSize  = 50

	SwordAttachmentWidth     = 165
	SwordAttachmentHeight    = 65
	SwordAttachmentSmoothing = 0.6

	PatheticSize     = 60
	PatheticMaxSpeed = 300
	PatheticHelth    = 100

	HaterSize             = 50
	HaterMaxSpeed         = 300
	HaterDistanse         = 450
	HaitersBolletDamage   = 10
	HaitersBolletSpeed    = 300
	HaitersBolletSize     = 15
	HaitersBolletCooldown = 1.0
)

var MaxPlayerSpeedMoving = 1000

var HeadAcceleration = 1000.0
var HeadNormalDeceleration = 800.0
var HeadDeceleration = HeadNormalDeceleration

var PatheticNormalDistanse = 450.0
var PatheticDistanse = 450.0
var PatheticNormalSpeed = 50.0
var PatheticBaseSpeed = PatheticNormalSpeed

var PathiticNormalDamage = 5
var PatheticDamage = PathiticNormalDamage

var NormalAcceleration = 500.0
var NormalDeceleration = 300.0

var Acceleration = NormalAcceleration
var Deceleration = NormalAcceleration

var PlayerDamage int
var PlayerHelth = MaxPlayerHelth
var MaxPlayerHelth = 100
var MaxSpeed = 300.0
var HeadRotationSpeed = 450.0
var Score int = 1

var MaxValwe int = 4
var Valwe int = 3

var MaxHatersValwe int = 5
var HatersValwe int = 2
var HaterNormalBaseSpeed = 50
var HaterBaseSpeed = HaterNormalBaseSpeed

type PointPlayer struct {
	Px, Py float64
}

type PointSpeed struct {
	Vx, Vy float64
}
