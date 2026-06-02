package enemy

import (
	"great-sword/game"
	"great-sword/game/common"
	"great-sword/game/player"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

var _ game.Entity = (*Haters)(nil)

type Haters struct {
	HatersMass []Heits
}

type Heits struct {
	HX, HY float64
	Active bool
}

func NewHaters() *Haters {
	h := &Haters{}
	return h
}

func (h *Haters) SpawnHiters() {
	x, y := RangomSpawnInWall()

	h.HatersMass = append(h.HatersMass, Heits{
		HX:     x,
		HY:     y,
		Active: true,
	})
}

func (h *Haters) UpdateHaiters(dt float64, world game.WorldView) {

	playerX, playerY, attahmentX, attahmentY, attahmentW, attahmentH, attahmentAngle := WhereThePlayer(world)

	for i := 0; i < len(h.HatersMass); i++ {
		enemy := &h.HatersMass[i]
		if !enemy.Active {
			continue
		}

		dx := playerX - enemy.HX
		dy := playerY - enemy.HY
		distance := math.Sqrt(dx*dx + dy*dy)

		var speed float64

		if distance > common.HaterDistanse {
			speed = common.HaterMaxSpeed
		} else if distance > common.HaterDistanse-20 && distance < common.HaterDistanse {
			speed = 0
		} else if distance < common.HaterDistanse {
			speed = -common.HaterBaseSpeed
		}
		if distance > 0.01 {
			enemy.HX += (dx / distance) * speed * dt
			enemy.HY += (dy / distance) * speed * dt
		}
		if CheckCollisionWithAttachment(enemy.HX, enemy.HY, common.HaterSize, attahmentX,
			attahmentY, float64(attahmentW), float64(attahmentH), attahmentAngle) {
			h.HatersMass = append(h.HatersMass[:i], h.HatersMass[i+1:]...)
			i--
			common.Score++

			player.ActivateBoost()

			continue
		}
		if CheckCollisionWithPlayer(enemy.HX, enemy.HY, common.HaterSize, world) {

			common.PlayerHelth -= common.HaitersBolletDamage
			h.HatersMass = append(h.HatersMass[:i], h.HatersMass[i+1:]...)
			i--

			//ИГРОК ДОЛЖЕН ОТТАЛКИВАТЬ ОТ СЕБЯ ВРАГОВ

		}

	}
}

func (h *Haters) Update(worldView game.WorldView) bool {
	dt := 1.0 / 60.0

	if len(h.HatersMass) < 2 {
		h.SpawnHiters()
	}

	h.UpdateHaiters(dt, worldView)
	return false
}

func (h *Haters) Draw(screen *ebiten.Image) {
}

func (h *Haters) Tag() string {
	return "hater"
}
