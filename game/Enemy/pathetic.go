package enemy

import (
	"great-sword/game"
	"great-sword/game/common"
	"great-sword/game/player"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

var _ game.Entity = (*Pathetic)(nil)

type Pathetic struct {
	Paths []OnePath
}

type OnePath struct {
	PX, PY float64
	Active bool
}

func NewPathetic() *Pathetic {
	p := &Pathetic{}
	return p
}

func (p *Pathetic) SpawnPathetic() {

	x, y := RangomSpawnInWall()

	p.Paths = append(p.Paths, OnePath{
		PX:     x,
		PY:     y,
		Active: true,
	})
}

func (p *Pathetic) UpatePathetics(dt float64, worldView game.WorldView) {

	playerX, playerY, attahmentX, attahmentY, attahmentW, attahmentH, attahmentAngle := WhereThePlayer(worldView)

	for i := 0; i < len(p.Paths); i++ {
		enemy := &p.Paths[i]
		if !enemy.Active {
			continue
		}

		dx := playerX - enemy.PX
		dy := playerY - enemy.PY
		distance := math.Sqrt(dx*dx + dy*dy)

		var speed float64

		if distance > common.PatheticDistanse {
			speed = common.PatheticBaseSpeed
		} else {
			t := 1.0 - distance/common.PatheticDistanse
			speed = common.PatheticBaseSpeed + t*float64(common.PatheticMaxSpeed-common.PatheticBaseSpeed)
		}
		if distance > 0.01 {
			enemy.PX += (dx / distance) * speed * dt
			enemy.PY += (dy / distance) * speed * dt
		}

		if CheckCollisionWithAttachment(enemy.PX, enemy.PY, common.PatheticSize, attahmentX,
			attahmentY, float64(attahmentW), float64(attahmentH), attahmentAngle) {
			p.Paths = append(p.Paths[:i], p.Paths[i+1:]...)
			i--

			if common.Score%10 == 0 {
				common.Valwe++
			}

			if common.Score >= 100 {
				if common.Score%50 == 0 {
					common.PatheticDamage++
				}
				if common.Score%10 == 0 {
					common.PatheticBaseSpeed += 10.0
				}
			}

			common.Score++

			if common.PlayerHelth <= common.MaxPlayerHelth-5 {
				common.PlayerHelth += 5
			}

			player.ActivateBoost()

			continue
		}
		if CheckCollisionWithPlayer(enemy.PX, enemy.PY, common.PatheticSize, worldView) {

			common.PlayerHelth -= common.PatheticDamage

			speed = -speed * 10

			enemy.PX += (dx / distance) * speed * dt
			enemy.PY += (dy / distance) * speed * dt

			//ИГРОК ДОЛЖЕН ОТТАЛКИВАТЬ ОТ СЕБЯ ВРАГОВ

		}

	}

}

func (p *Pathetic) Update(worldView game.WorldView) bool {
	dt := 1.0 / 60.0

	if common.Valwe > common.MaxValwe {
		common.Valwe = common.MaxValwe
	}

	number := common.Valwe

	if len(p.Paths) < number {
		p.SpawnPathetic()
	}

	p.UpatePathetics(dt, worldView)

	return false
}

func (p *Pathetic) Draw(screen *ebiten.Image) {
}

func (p *Pathetic) Tag() string {
	return "pathetic"
}
