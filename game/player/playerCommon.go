package player

import "great-sword/game/common"

var SwordIxist bool = true
var SwordIxistTimer float64
var SwordIxistTimerNormal float64 = 0.3

func SwordVanished() {
	if SwordIxistTimer <= 0 {
		SwordIxistTimer = SwordIxistTimerNormal
		if SwordIxist {
			SwordIxist = false
		} else {
			SwordIxist = true
		}
	}
}

var BoostTimer float64
var NormalSpeed float64 = common.MaxSpeed
var NormalRotSpeed float64 = common.HeadRotationSpeed
var NormalBoostTimerLong float64 = 0.5
var BoostTimerLong float64 = NormalBoostTimerLong

var NormalBoostRotating int = 300
var NormalBoostSpeed int = 300

var MaxBoostSpeed int = 450

var BoostRotating int = NormalBoostRotating
var BoostSpeed int = NormalBoostSpeed

func ActivateBoost() {
	if BoostTimer <= 0 {
		BoostTimer = BoostTimerLong
		BoostRotating = NormalBoostRotating
		BoostSpeed = NormalBoostSpeed
		BoostTimerLong = NormalBoostTimerLong
	} else if BoostRotating <= MaxBoostSpeed-20 && BoostSpeed <= MaxBoostSpeed-20 {
		BoostRotating += 20
		BoostSpeed += 20
		BoostTimerLong += 0.2
	}
}

var Forward int = 1000
var ForwardTimer float64
var ForwadTimerLong float64 = 0.5
var RechargeForwartTimer float64
var RechargeForwartTimerLong = 2.0

func forward() {
	ForwardTimer = ForwadTimerLong
	common.Deceleration = 0
	common.Acceleration = common.Acceleration * 2
}

func ActivatedForward() {
	RechargeForwartTimer = RechargeForwartTimerLong
	forward()
}

var Rebount float64 = 0.6
