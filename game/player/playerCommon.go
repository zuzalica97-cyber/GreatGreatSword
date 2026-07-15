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
var NormalBoostTimerLong float64 = 0.5
var BoostTimerLong float64 = NormalBoostTimerLong

var NormalBoostRotating int = 150
var NormalBoostSpeed int = 150

var MaxBoostSpeed int = 350

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

var Rebount float64 = 0.6

func (p *PlayerLeg) DamagePlayer(damage int) {
	if p.Speed.Vx < 600 && p.Speed.Vx >= -600 && p.Speed.Vy <= 600 && p.Speed.Vy >= -600 {
		common.PlayerHelth -= float64(damage)
	}
}
