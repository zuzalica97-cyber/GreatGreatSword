package player

import (
	"fmt"
	"great-sword/game"
	"great-sword/game/common"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

var _ game.Entity = (*PlayerLeg)(nil)

type PlayerLeg struct {
	Position common.PointPlayer
	Speed    common.PointSpeed
}

func NewPlayerLeg() *PlayerLeg {

	return &PlayerLeg{
		Position: common.PointPlayer{
			Px: common.ScreenWidth/2 - common.PlayerSize/2,
			Py: common.ScreenHeight/2 - common.PlayerSize/2,
		},
		Speed: common.PointSpeed{
			Vx: 0,
			Vy: 0,
		},
	}
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func (p *PlayerLeg) ActivateBoost() {
	BoostActive = true
	BoostTimer = BoostTimerLong
}

func (p *PlayerLeg) Update(worldView game.WorldView) bool {

	dt := 1.0 / 60.0

	if BoostActive {
		BoostTimer -= dt
		if BoostTimer <= 0 {
			BoostActive = false
		}
	}

	if ForwardActive {
		ForwardTimer -= dt
		if ForwardTimer <= 0 {
			ForwardActive = false
			common.Deceleration = common.NormalDeceleration
			common.Acceleration = common.NormalAcceleration
		}
	}

	if RechargeForwart {
		RechargeForwartTimer -= dt
		if RechargeForwartTimer <= 0 {
			RechargeForwart = false
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if !RechargeForwart {
			ActivatedForward()
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyT) {
		fmt.Println(p.Speed.Vx, p.Speed.Vy)
	}

	rebound := Rebount

	axelerationLeg := common.Acceleration

	currentMaxSpeed := common.MaxSpeed

	if BoostActive {
		currentMaxSpeed = common.MaxSpeed + float64(BoostSpeed)
	}
	if ForwardActive {
		rebound = Rebount / 2
		axelerationLeg = common.Acceleration + 500.0
		currentMaxSpeed = common.MaxSpeed + float64(Forward) // НУЖНО УВЕЛИЧЕТЬ СКОРСТЬ ПРИ ОБЫЧНОМ ИСПОЛЬЗОВАНИИ НЕ УВЕЛИЧЕВАЯ  ИМЕЮЩИЙСЯ ОТСОК от стен
	}
	if ForwardActive && BoostActive {
		currentMaxSpeed = common.MaxSpeed + float64(Forward) + float64(BoostSpeed)*0.5
	}

	if currentMaxSpeed > float64(common.MaxPlayerSpeedMoving) {
		currentMaxSpeed = float64(common.MaxPlayerSpeedMoving)
	}

	moveX, moveY := 0.0, 0.0

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		moveX = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		moveX = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		moveY = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		moveY = 1
	}

	if moveX != 0 && moveY != 0 { //движение подиоганале
		moveX *= 0.7071
		moveY *= 0.7071
	}

	if moveX != 0 { //Применяем ускорение
		p.Speed.Vx += moveX * axelerationLeg * dt
	} else {
		dec := common.Deceleration * dt // Замедление если ненажата клавиша
		if math.Abs(p.Speed.Vx) > dec {
			p.Speed.Vx -= math.Copysign(dec, p.Speed.Vx)
		} else {
			p.Speed.Vx = 0
		}
	}

	if moveY != 0 { //Применяем ускорение
		p.Speed.Vy += moveY * axelerationLeg * dt
	} else {
		dec := common.Deceleration * dt // Замедление если ненажата клавиша
		if math.Abs(p.Speed.Vy) > dec {
			p.Speed.Vy -= math.Copysign(dec, p.Speed.Vy)
		} else {
			p.Speed.Vy = 0
		}
	}

	if p.Speed.Vx > currentMaxSpeed {
		p.Speed.Vx = currentMaxSpeed
	}
	if p.Speed.Vx < -currentMaxSpeed {
		p.Speed.Vx = -currentMaxSpeed
	}
	if p.Speed.Vy > currentMaxSpeed {
		p.Speed.Vy = currentMaxSpeed
	}
	if p.Speed.Vy < -currentMaxSpeed {
		p.Speed.Vy = -currentMaxSpeed
	}

	p.Position.Px += p.Speed.Vx * dt
	p.Position.Py += p.Speed.Vy * dt

	//Границы с отскоком и потерей скорости
	if p.Position.Px < 0 {
		p.Position.Px = 0
		p.Speed.Vx = -p.Speed.Vx * rebound //отскок с потерей скорости
	}
	if p.Position.Px > common.ScreenWidth-common.PlayerSize {
		p.Position.Px = common.ScreenWidth - common.PlayerSize
		p.Speed.Vx = -p.Speed.Vx * rebound //отскок с потерей скорости
	}
	if p.Position.Py < 0 {
		p.Position.Py = 0
		p.Speed.Vy = -p.Speed.Vy * rebound //отскок с потерей скорости
	}
	if p.Position.Py > common.ScreenHeight-common.PlayerSize {
		p.Position.Py = common.ScreenHeight - common.PlayerSize
		p.Speed.Vy = -p.Speed.Vy * rebound //отскок с потерей скорости
	}
	return false
}

func (p *PlayerLeg) Draw(screen *ebiten.Image) {
}

func (p *PlayerLeg) Tag() string {
	return "playerLeg"
}
