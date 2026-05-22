package swords

import (
	"great-sword/game"
	"great-sword/game/common"
)

var _ game.Entity = (*BlueSword)(nil)

type BlueSword struct {
	Position         common.PointPlayer
	TargetX, TargetY float64
	Angle            float64
	TargetAngle      float64
}
