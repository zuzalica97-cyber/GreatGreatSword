package enemy

import (
	"fmt"
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
	PX, PY                 float64
	Helth                  int
	PathericCooldownActive bool
	PathericCooldownTimer  float64
	Active                 bool
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
		Helth:  common.PatheticHelth,
	})
}

func (p *Pathetic) UpatePathetics(dt float64, worldView game.WorldView) {

	playerX, playerY, attahmentX, attahmentY, attahmentW, attahmentH, attahmentAngle := WhereThePlayer(worldView)

	for i := 0; i < len(p.Paths); i++ {
		enemy := &p.Paths[i]
		if !enemy.Active {
			continue
		}

		if enemy.PathericCooldownActive {
			enemy.PathericCooldownTimer -= dt
			if enemy.PathericCooldownTimer <= 0 {
				enemy.PathericCooldownActive = false
			}
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
		if enemy.PathericCooldownActive {
			speed = -100
		}

		if distance > 0.01 {
			enemy.PX += (dx / distance) * speed * dt
			enemy.PY += (dy / distance) * speed * dt
		}

		if CheckCollisionWithAttachment(enemy.PX, enemy.PY, common.PatheticSize, attahmentX,
			attahmentY, float64(attahmentW), float64(attahmentH), attahmentAngle) {

			if !enemy.PathericCooldownActive || enemy.PathericCooldownTimer <= 0.5 {
				enemy.Helth -= common.PlayerDamage
				fmt.Println(common.PlayerDamage)
				PatheticCooldown(enemy)
			}

			if common.Score%20 == 0 {
				common.Valwe++
			}

			if common.Score >= 100 {
				if common.Score%100 == 0 {
					common.PatheticDamage++
				}
				if common.Score%10 == 0 {
					common.PatheticBaseSpeed += 10.0
				}
			}

			if common.Score%20 == 0 {
				common.Valwe++
			}

			if common.Score >= 100 {
				if common.Score%100 == 0 {
					common.PatheticDamage++
				}
				if common.Score%50 == 0 {
					common.PatheticBaseSpeed += 10.0
				}
				if common.Score%20 == 0 {
					common.PatheticDistanse += 10
				}
			}

			if enemy.Helth <= 0 {
				p.Paths = append(p.Paths[:i], p.Paths[i+1:]...)
				i--
				common.Score++
				common.PlayerHelth += 2
				player.ActivateBoost()
			}

			continue
		}
		if CheckCollisionWithPlayer(enemy.PX, enemy.PY, common.PatheticSize, worldView) {

			if !enemy.PathericCooldownActive {
				common.PlayerHelth -= common.PatheticDamage
			}

			PatheticCooldown(enemy)

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
