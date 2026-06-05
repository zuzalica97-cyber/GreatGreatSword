package player

import (
	"great-sword/game"
	"great-sword/game/common"

	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

var _ game.Entity = (*PlayerHead)(nil)

type PlayerHead struct {
	PositionHead    common.PointPlayer
	Angle           float64
	AngularVelocity float64

	BoostTimer     float64
	BoostActive    bool
	NormalSpeed    float64
	NormalROTSPEED float64
}

func NewPlayerHead() *PlayerHead {
	return &PlayerHead{
		Angle:           0.0,
		AngularVelocity: 0.0,
	}
}

func (p *PlayerHead) Update(woroldWiev game.WorldView) bool {
	dt := 1.0 / 60.0

	currentRotSpeed := common.HeadRotationSpeed

	if BoostActive {
		currentRotSpeed = common.MaxSpeed + float64(BoostRotating)
	}
	if ForwardActive && BoostActive {
		currentRotSpeed = common.MaxSpeed + float64(BoostRotating)*2.5
	}

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

	if p.AngularVelocity > currentRotSpeed {
		p.AngularVelocity = currentRotSpeed
	}
	if p.AngularVelocity < -currentRotSpeed {
		p.AngularVelocity = -currentRotSpeed
	}

	p.Angle += p.AngularVelocity * dt

	p.Angle = math.Mod(p.Angle, 360)
	if p.Angle < 0 {
		p.Angle += 360
	}

	for _, entity := range woroldWiev.SearchEntities("playerLeg") {
		pL := entity.(*PlayerLeg)

		centerX := pL.Position.Px + common.PlayerSize/2
		centerY := pL.Position.Py + common.PlayerSize/2

		p.PositionHead.Px = centerX
		p.PositionHead.Py = centerY

		vx := pL.Speed.Vx
		vy := pL.Speed.Vy

		if vx < 0 {
			vx = -vx
		}
		if vy < 0 {
			vy = -vy
		}
		if currentRotSpeed < 0 {
			currentRotSpeed = -currentRotSpeed
		}

		FinalDamage := ((pL.Speed.Vx + pL.Speed.Vy) * 0.1) + (currentRotSpeed * 0.1) + 50

		if p.BoostActive {
			FinalDamage = FinalDamage * 1.5
		}

		common.PlayerDamage = int(FinalDamage)
	}

	return false
}

func (p *PlayerHead) Draw(screen *ebiten.Image) {

}

func (p *PlayerHead) Tag() string {
	return "playerHead"
}
