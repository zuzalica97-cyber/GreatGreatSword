package player

import (
	"great-sword/game"
	"great-sword/game/common"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

var _ game.Entity = (*PlayerHead)(nil)

type PlayerHead struct {
	Angle           float64
	AngularVelocity float64
}

func NewPlayerHead() *PlayerHead {
	return &PlayerHead{
		Angle:           0.0,
		AngularVelocity: 0.0,
	}
}

func (p *PlayerHead) Update(woroldWiev game.WorldView) bool {
	dt := 1.0 / 60.0

	rotationInput := 0.0

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		rotationInput = -1
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		rotationInput = 1
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) &&
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		rotationInput = 0
	}

	if rotationInput != 0 {
		p.AngularVelocity += rotationInput * common.HeadAcceleration * dt
	} else {
		dec := common.HeadDeceleration * dt
		if math.Abs(p.AngularVelocity) > dec {
			p.AngularVelocity -= math.Copysign(dec, p.AngularVelocity)
		} else {
			p.AngularVelocity = 0
		}
	}

	p.AngularVelocity = clamp(p.AngularVelocity, -common.HeadRotationSpeed, common.HeadRotationSpeed)

	p.Angle += p.AngularVelocity * dt

	p.Angle = math.Mod(p.Angle, 360)
	if p.Angle < 0 {
		p.Angle += 360
	}

	return false
}

func (p *PlayerHead) Draw(screen *ebiten.Image) {

}

func (p *PlayerHead) Tag() string {
	return "playerHead"
}
