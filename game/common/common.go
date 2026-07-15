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

	HaterSize             = 50.0
	HaterMaxSpeed         = 300
	HaterDistanse         = 600
	HaitersBolletDamage   = 10
	HaitersBolletSpeed    = 300
	HaitersBolletSize     = 15.0
	HaitersBolletCooldown = 1.0
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

var MaxValwe int = 4
var Valwe int = 3

var MaxHatersValwe int = 5
var HatersValwe int = 2
var HaterNormalBaseSpeed = 150
var HaterBaseSpeed = HaterNormalBaseSpeed

type PointPlayer struct {
	Px, Py float64
}

type PointSpeed struct {
	Vx, Vy float64
}
