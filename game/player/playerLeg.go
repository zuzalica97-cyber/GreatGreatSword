package player

import (
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

func (p *PlayerLeg) ResetGame(world game.WorldView) {
	common.Score = 0

	p.Position.Px = common.ScreenWidth/2 - common.PlayerSize/2
	p.Position.Py = common.ScreenHeight/2 - common.PlayerSize/2
	p.Speed.Vx = 0
	p.Speed.Vy = 0

	for _, entity := range world.SearchEntities("playerHead") {
		head := entity.(*PlayerHead)

		head.Angle = 0
		head.AngularVelocity = 0

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

func (p *PlayerLeg) Update(worldView game.WorldView) bool {

	dt := 1.0 / 60.0

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
		p.Speed.Vx += moveX * common.Acceleration * dt
	} else {
		dec := common.Deceleration * dt // Замедление если ненажата клавиша
		if math.Abs(p.Speed.Vx) > dec {
			p.Speed.Vx -= math.Copysign(dec, p.Speed.Vx)
		} else {
			p.Speed.Vx = 0
		}
	}

	if moveY != 0 { //Применяем ускорение
		p.Speed.Vy += moveY * common.Acceleration * dt
	} else {
		dec := common.Deceleration * dt // Замедление если ненажата клавиша
		if math.Abs(p.Speed.Vy) > dec {
			p.Speed.Vy -= math.Copysign(dec, p.Speed.Vy)
		} else {
			p.Speed.Vy = 0
		}
	}

	// Ограничение максимальной скорости
	p.Speed.Vx = clamp(p.Speed.Vx, -common.MaxSpeed, common.MaxSpeed)
	p.Speed.Vy = clamp(p.Speed.Vy, -common.MaxSpeed, common.MaxSpeed)

	p.Position.Px += p.Speed.Vx * dt
	p.Position.Py += p.Speed.Vy * dt

	//Границы с отскоком и потерей скорости
	if p.Position.Px < 0 {
		p.Position.Px = 0
		p.Speed.Vx = -p.Speed.Vx * 0.6 //отскок с потерей скорости
	}
	if p.Position.Px > common.ScreenWidth-common.PlayerSize {
		p.Position.Px = common.ScreenWidth - common.PlayerSize
		p.Speed.Vx = -p.Speed.Vx * 0.6 //отскок с потерей скорости
	}
	if p.Position.Py < 0 {
		p.Position.Py = 0
		p.Speed.Vy = -p.Speed.Vy * 0.6 //отскок с потерей скорости
	}
	if p.Position.Py > common.ScreenHeight-common.PlayerSize {
		p.Position.Py = common.ScreenHeight - common.PlayerSize
		p.Speed.Vy = -p.Speed.Vy * 0.6 //отскок с потерей скорости
	}
	return false
}

func (p *PlayerLeg) Draw(screen *ebiten.Image) {
}

func (p *PlayerLeg) Tag() string {
	return "playerLeg"
}
