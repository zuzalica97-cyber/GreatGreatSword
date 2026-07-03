package player

import (
	"great-sword/game"
	"great-sword/game/common"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
)

var _ game.Entity = (*PlayerLeg)(nil)

type PlayerLeg struct {
	Position common.PointPlayer
	Speed    common.PointSpeed
	Texture  *ebiten.Image
}

func NewPlayerLeg() *PlayerLeg {

	return &PlayerLeg{
		Position: common.PointPlayer{
			Px: common.RoomWidth/2 - common.PlayerSize/2,
			Py: common.RoomHeight/2 - common.PlayerSize/2,
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
	BoostTimer = BoostTimerLong
}

func (p *PlayerLeg) Update(worldView game.WorldView) bool {

	dt := 1.0 / 60.0

	if common.PlayerHelth >= common.MaxPlayerHelth {
		common.PlayerHelth = common.MaxPlayerHelth
	}

	if BoostTimer > 0 {
		BoostTimer -= dt
	}

	if SwordIxistTimer > 0 {
		SwordIxistTimer -= dt
	}

	if ForwardTimer > 0 {
		ForwardTimer -= dt
		if ForwardTimer <= 0 {
			common.Deceleration = common.NormalDeceleration
			common.Acceleration = common.NormalAcceleration
		}
	}

	if RechargeForwartTimer > 0 {
		RechargeForwartTimer -= dt
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if RechargeForwartTimer <= 0 {
			ActivatedForward()
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyF) {
		SwordVanished()
	}

	rebound := Rebount

	axelerationLeg := common.Acceleration
	deAxeleration := common.Deceleration

	currentMaxSpeed := common.MaxSpeed

	if BoostTimer > 0 {
		currentMaxSpeed = common.MaxSpeed + float64(BoostSpeed)
	}
	if ForwardTimer > 0 {
		rebound = Rebount / 2
		axelerationLeg = common.Acceleration + 500.0
		currentMaxSpeed = common.MaxSpeed + float64(Forward)*10 // НУЖНО УВЕЛИЧЕТЬ СКОРСТЬ ПРИ ОБЫЧНОМ ИСПОЛЬЗОВАНИИ НЕ УВЕЛИЧЕВАЯ  ИМЕЮЩИЙСЯ ОТСОК от стен
	}
	if ForwardTimer > 0 && ForwardTimer > 0 {
		currentMaxSpeed = common.MaxSpeed + float64(Forward) + float64(BoostSpeed)*10
	}
	if SwordIxistTimer > 0 {
		deAxeleration = common.Deceleration / 10
		axelerationLeg = common.Acceleration / 5
		speed := currentMaxSpeed
		currentMaxSpeed = speed * 2
		forwardT := ForwadTimerLong
		rechangeForwT := RechargeForwartTimerLong
		ForwadTimerLong = forwardT * 2
		RechargeForwartTimerLong = rechangeForwT / 2
	} else {
		ForwadTimerLong = 0.5
		RechargeForwartTimerLong = 2.0
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
		dec := deAxeleration * dt // Замедление если ненажата клавиша
		if math.Abs(p.Speed.Vx) > dec {
			p.Speed.Vx -= math.Copysign(dec, p.Speed.Vx)
		} else {
			p.Speed.Vx = 0
		}
	}

	if moveY != 0 { //Применяем ускорение
		p.Speed.Vy += moveY * axelerationLeg * dt
	} else {
		dec := deAxeleration * dt // Замедление если ненажата клавиша
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
	if p.Position.Px > common.RoomWidth-common.PlayerSize {
		p.Position.Px = common.RoomHeight - common.PlayerSize
		p.Speed.Vx = -p.Speed.Vx * rebound //отскок с потерей скорости
	}
	if p.Position.Py < 0 {
		p.Position.Py = 0
		p.Speed.Vy = -p.Speed.Vy * rebound //отскок с потерей скорости
	}
	if p.Position.Py > common.RoomHeight-common.PlayerSize {
		p.Position.Py = common.RoomHeight - common.PlayerSize
		p.Speed.Vy = -p.Speed.Vy * rebound //отскок с потерей скорости
	}

	return false
}

func (p *PlayerLeg) Draw(screen *ebiten.Image, camera *kamera.Camera) {
}

func (p *PlayerLeg) Tag() string {
	return "playerLeg"
}
